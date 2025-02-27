package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger
var rsn string

func main() {
	logger = log.Default()

	rsn = os.Getenv("PLAYER_NAME")
	if rsn == "" {
		logger.Fatal("PLAYER_NAME environment variable not defined, exiting...")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8340"
	}

	logger.Printf("Web server started on port %s!\n", port)
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", metrics)	
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		logger.Fatalf("Error creating the web server: %s", err)
	}
}

func metrics(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://api.wiseoldman.net/v2/players/" + rsn)
	if err != nil {
		w.Write([]byte(""))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		w.Write([]byte(""))
		return
	}
	
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		w.Write([]byte(""))
	}

	output := ""

	output += fmt.Sprintf("runescape_combatlevel{player=\"%s\"} %.0f\n", rsn, data["combatLevel"])

	for _, skill := range data["latestSnapshot"].(map[string]interface{})["data"].(map[string]interface{})["skills"].(map[string]interface{}) {
		if skillMap, ok := skill.(map[string]interface{}); ok {
			output += fmt.Sprintf("runescape_%s_level{player=\"%s\"} %.0f\n", skillMap["metric"], rsn, skillMap["level"])
			output += fmt.Sprintf("runescape_%s_experience{player=\"%s\"} %.0f\n", skillMap["metric"], rsn, skillMap["experience"])
			output += fmt.Sprintf("runescape_%s_rank{player=\"%s\"} %.0f\n", skillMap["metric"], rsn, skillMap["rank"])
			output += fmt.Sprintf("runescape_%s_ehp{player=\"%s\"} %f\n", skillMap["metric"], rsn, skillMap["ehp"])
		}
	}

	for _, boss := range data["latestSnapshot"].(map[string]interface{})["data"].(map[string]interface{})["bosses"].(map[string]interface{}) {
		if bossMap, ok := boss.(map[string]interface{}); ok {
			output += fmt.Sprintf("runescape_%s_kills{player=\"%s\"} %.0f\n", bossMap["metric"], rsn, bossMap["kills"])
			output += fmt.Sprintf("runescape_%s_rank{player=\"%s\"} %.0f\n", bossMap["metric"], rsn, bossMap["rank"])
			output += fmt.Sprintf("runescape_%s_ehb{player=\"%s\"} %f\n", bossMap["metric"], rsn, bossMap["ehb"])
		}
	}

	for _, activity := range data["latestSnapshot"].(map[string]interface{})["data"].(map[string]interface{})["activities"].(map[string]interface{}) {
		if activityMap, ok := activity.(map[string]interface{}); ok {
			output += fmt.Sprintf("runescape_%s_score{player=\"%s\"} %.0f\n", activityMap["metric"], rsn, activityMap["score"])
			output += fmt.Sprintf("runescape_%s_rank{player=\"%s\"} %.0f\n", activityMap["metric"], rsn, activityMap["rank"])
		}
	}

	for _, total := range data["latestSnapshot"].(map[string]interface{})["data"].(map[string]interface{})["computed"].(map[string]interface{}) {
		if totalMap, ok := total.(map[string]interface{}); ok {
			output += fmt.Sprintf("runescape_%s{player=\"%s\"} %f\n", totalMap["metric"], rsn, totalMap["value"])
			output += fmt.Sprintf("runescape_%s_rank{player=\"%s\"} %.0f\n", totalMap["metric"], rsn, totalMap["rank"])
		}
	}

	w.Write([]byte(output))
}