package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/battlesnakeio/exporter/engine"
	"github.com/battlesnakeio/exporter/render"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "https://battlesnake.io", 302)
}

func handleASCIIFrame(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	gameID := p.ByName("game")
	frameID, err := strconv.Atoi(p.ByName("frame"))
	if err != nil {
		handleError(w, r, err)
		return
	}

	game, err := engine.GetGame(gameID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	gameFrame, err := engine.GetGameFrame(game.ID, frameID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	if err = render.GameFrameToASCII(w, game, gameFrame); err != nil {
		handleError(w, r, err)
		return
	}
}

func handleGIFFrame(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	gameID := p.ByName("game")
	frameID, err := strconv.Atoi(p.ByName("frame"))
	if err != nil {
		handleError(w, r, err)
		return
	}

	game, err := engine.GetGame(gameID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	gameFrame, err := engine.GetGameFrame(game.ID, frameID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "image/gif")
	if err = render.GameFrameToGIF(w, game, gameFrame); err != nil {
		handleError(w, r, err)
		return
	}
}

func handleGIFGame(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	gameID := p.ByName("game")
	game, err := engine.GetGame(gameID)
	if err != nil {
		handleError(w, r, err)
		return
	}

	frames := strings.Split(r.URL.Query().Get("frames"), "-")
	offset := 0
	frameRange := -1
	if len(frames) == 2 {
		var err error
		var endFrame int
		offset, err = strconv.Atoi(frames[0])
		if err != nil {
			log.WithError(err).Errorf("unable to convert offset: %s", frames[0])
			offset = 0
		}
		endFrame, err = strconv.Atoi(frames[1])
		if err != nil {
			log.WithError(err).Errorf("unable to convert ending frame: %s", frames[1])
		} else {
			frameRange = endFrame - offset + 1
		}
	}

	gameFrames, err := engine.GetGameFrames(game.ID, offset, frameRange)
	if err != nil {
		handleError(w, r, err)
		return
	}

	frameDelay, err := strconv.Atoi(r.URL.Query().Get("frameDelay"))
	if err != nil {
		frameDelay = render.GIFFrameDelay
	}

	loopDelay, err := strconv.Atoi(r.URL.Query().Get("loopDelay"))
	if err != nil {
		loopDelay = render.GIFLoopDelay
	}

	w.Header().Set("Content-Type", "image/gif")
	err = render.GameFramesToAnimatedGIF(w, game, gameFrames, frameDelay, loopDelay)
	if err != nil {
		handleError(w, r, err)
		return
	}
}

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	log.WithError(err).
		WithFields(log.Fields{
			"httpRequest": map[string]interface{}{
				"method":    r.Method,
				"url":       r.URL.String(),
				"userAgent": r.Header.Get("User-Agent"),
				"referrer":  r.Header.Get("Referer"),
			},
		}).Error("unable to process request")
	w.WriteHeader(http.StatusInternalServerError)
	if _, err := w.Write([]byte(err.Error())); err != nil {
		log.WithError(err).Error("unable to write to response stream")
	}
	return
}
