package main

// [START import]
import (
        "context"
        "fmt"
        "log"
        "net/http"
        "os"
        "strconv"

        "cloud.google.com/go/datastore"
)

// [END import]
// [START main_func]
type Counter struct {
        Count int
}

func main() {
        http.HandleFunc("/", indexHandler)

        //    [START setting_port]
        port := os.Getenv("PORT")
        if port == "" {
                port = "8080"
                log.Printf("Defaulting to port %s", port)
        }

        log.Printf("Listening on port %s", port)
        if err := http.ListenAndServe(":"+port, nil); err != nil {
                log.Fatal(err)
        }
        // [END setting_port]
}

// [END main_func]

// [START indexHandler]

// indexHandler responds to requests with our greeting.
func indexHandler(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/version" {
                fmt.Fprintf(w, "This is app verison B!")
        }
        if r.URL.Path == "/instance" {
                fmt.Fprintf(w, os.Getenv("GAE_INSTANCE")+" "+os.Getenv("DEVSHELL_PROJECT_ID"))
        }
        if r.URL.Path == "/" {
                ctx := context.Background()
                projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
                client, err := datastore.NewClient(ctx, projectID)
                if err != nil {
                        http.Error(w, "Failed to create client: ", http.StatusOK)
                }
                kind := "Counter"
                name := "count1"
                taskKey := datastore.NameKey(kind, name, nil)
                curCount := Counter{
                        Count: 1,
                }
                //If the client doesnt retrieve the counter struct then it adds the the created struct.
                //else it overwrites the current struct and increaments its value.
                if err1 := client.Get(ctx, taskKey, &curCount); err1 != nil {
                        if _, err2 := client.Put(ctx, taskKey, &curCount); err2 != nil {
                                strE := fmt.Sprintf("%s", err2)
                                http.Error(w, strE, http.StatusOK)
                        }
                } else {
                        curCount.Count += 1
                        if _, err3 := client.Put(ctx, taskKey, &curCount); err3 != nil {
                                strE3 := fmt.Sprintf("%s", err3)
                                http.Error(w, strE3, http.StatusOK)
                        }
                }
                w.Write([]byte("Joey Bennie 300280586 visit counter: " + strconv.Itoa(curCount.Count)))
                defer client.Close()
        }
}
