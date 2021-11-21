package nanny

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func (n *Nanny) timeRemaining(w http.ResponseWriter, r *http.Request) {
	type templateData struct {
		RemainingTime string
	}
	availTimeDur := time.Duration(n.state.AvailableTimeSec) * time.Second
	td := templateData{
		RemainingTime: fmt.Sprintf("%s", availTimeDur),
	}
	const tpl = `
<!DOCTYPE html>
<html>
<head>
	<title>Time remaining</title>
</head>
<body>
	<h2>GO-NANNY</h2>
	<h3>Time remaining</h3>
	<p>{{ .RemainingTime }}</p>
</body>
</html>
	`
	t, err := template.New("time-available").Parse(tpl)
	if err != nil {
		log.Error(err)
	}
	err = t.Execute(w, td)
	if err != nil {
		log.Error(err)
	}
}

func (n *Nanny) runServer() {
	r := mux.NewRouter()
	r.HandleFunc("/time-remaining", n.timeRemaining)

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8544",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
