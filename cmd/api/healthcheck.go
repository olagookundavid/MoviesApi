package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// data := map[string]string{
	// 	"status":      "available",
	// 	"environment": app.config.env,
	// 	"version":     version,
	// }
	env := envelope{"status": "available",
		"system_info": map[string]string{"environment": app.config.env,
			"version": version}}
	//err := json.NewEncoder(w).Encode(data)
	//this methods writes to output stream in a single step
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

}
