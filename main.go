package main

import (
	"secrethitler.io/database"
	"secrethitler.io/routes"
	//	"secrethitler.io/types"

	"context"
	"crypto/rand"
	"encoding/hex"
	//	"encoding/json"
	"fmt"
	"net/http"
	// "os"
	//	"strconv"
	"sort"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/googollee/go-socket.io"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/discord"
	"github.com/markbates/goth/providers/github"
	"github.com/spf13/viper"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	viper.SetDefault("ENV", "development")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("HOST", "")
	viper.SetDefault("COOKIE_MAX_AGE", 86400 * 30)
	viper.SetDefault("MONGODB_HOST", "localhost")
	viper.SetDefault("MONGODB_PORT", "27017")
	viper.SetDefault("MONGODB_NAME", "secret-hitler-app")
	viper.SetDefault("REDIS_HOST", "127.0.0.1")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_ID", 0)
	env, _ := viper.Get("ENV").(string)
	host, _ := viper.Get("HOST").(string)
	port, _ := viper.Get("PORT").(string)
	cacheToken, _ := viper.Get("CACHE_TOKEN").(string)
	sessionKey, _ := viper.Get("SESSION_KEY").(string)
	cookieMaxAge, _ := viper.Get("COOKIE_MAX_AGE").(int)
	mongoDBHost, _ := viper.Get("MONGODB_HOST").(string)
	mongoDBPort, _ := viper.Get("MONGODB_PORT").(string)
	mongoDBName, _ := viper.Get("MONGODB_NAME").(string)
	redisHost, _ := viper.Get("REDIS_HOST").(string)
	redisPort, _ := viper.Get("REDIS_PORT").(string)
	redisPass, _ := viper.Get("REDIS_PASS").(string)
	redisID, _ := viper.Get("REDIS_ID").(int)
	oauthRedirectHost, _ := viper.Get("OAUTH_REDIRECT_HOST").(string)
	discordKey, _ := viper.Get("DISCORD_KEY").(string)
	discordSecret, _ := viper.Get("DISCORD_SECRET").(string)
	githubKey, _ := viper.Get("GITHUB_KEY").(string)
	githubSecret, _ := viper.Get("GITHUB_SECRET").(string)
	writeConfig, _ := viper.Get("WRITE_CONFIG").(bool)

	if cacheToken == "" {
		bytes := make([]byte, 4)
		rand.Read(bytes)
		cacheToken = hex.EncodeToString(bytes)
	}

	if sessionKey == "" {
		bytes := make([]byte, 32)
		rand.Read(bytes)
		fmt.Println("Empty session key detected: set SESSION_KEY="+hex.EncodeToString(bytes)+" in .env.")
	}

	if writeConfig {
		 viper.WriteConfig()
	}

	routes.CacheToken = cacheToken
	fmt.Println(env, routes.CacheToken)
	uri := "mongodb://" + mongoDBHost + ":" + mongoDBPort + "/" + mongoDBName
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	defer cancel()
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	mongoDB := client.Database(mongoDBName)

	redisDB := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPass,
		DB:       redisID,
	})

	database.SetupDatabase(mongoDB, redisDB)

	goth.UseProviders(
		discord.New(discordKey, discordSecret, oauthRedirectHost+"/auth/discord/callback"),
		github.New(githubKey, githubSecret, oauthRedirectHost+"/auth/github/callback"),
	)

	oauthProviderMap := make(map[string]string)
	oauthProviderMap["discord"] = "Discord"
	oauthProviderMap["github"] = "Github"
	var oauthProviders []string

	for oauthProvider := range oauthProviderMap {
		oauthProviders = append(oauthProviders, oauthProvider)
	}

	sort.Strings(oauthProviders)
	io := socketio.NewServer(nil)

	store := sessions.NewCookieStore([]byte(sessionKey))
	store.MaxAge(cookieMaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = env == "production"

	router := mux.NewRouter()
	routes.SetupRoutes(router, io, store)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./public")))

	go io.Serve()
	defer io.Close()
	fmt.Println("Listening on " + host + ":" + port)
	http.ListenAndServe(host+":"+port, router)
}
