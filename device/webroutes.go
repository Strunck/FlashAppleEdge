package device

import (
	"fmt"
	"net/http"

	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
)

const LatestMessageCount = 30

func HandleLogMsg(ds State) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sse := datastar.NewSSE(w, r)

		msgBufferChannel := make(chan *nats.Msg, 30)
		nsub, err := ds.chanLatestMsg(msgBufferChannel)
		if err != nil {
			http.Error(w, "error subscribing to NATS messages", http.StatusInternalServerError)
			return
		}
		defer nsub.Unsubscribe()

		// Get latest NATS messages, until r.Context() is done
		for {
			select {
			case <-r.Context().Done():
				return
			case msg := <-msgBufferChannel:
				strBuffer := []LatestNatsMessage{
					{
						Subject: msg.Subject,
						Data:    string(msg.Data),
					},
				}

			drainLoop:
				for len(strBuffer) < LatestMessageCount {
					select {
					case additionalMsg := <-msgBufferChannel:
						strBuffer = append(strBuffer, LatestNatsMessage{
							Subject: additionalMsg.Subject,
							Data:    string(additionalMsg.Data),
						})
					default:
						break drainLoop
					}
				}

				if err := sse.PatchElementTempl(LatestNatsMessages(strBuffer)); err != nil {
					http.Error(w, "Error initializing DataStar SSE", http.StatusInternalServerError)
					return
				}
			}
		}
	})
}

func HandleIndexTbl(ds State) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sse := datastar.NewSSE(w, r)

		err := sse.PatchElements(loadIndexTbl(ds))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})
}

func loadIndexTbl(ds State) string {
	html := ""
	for l, line := range ds.Cfg.Lines {
		for i := 0; i < line.IdCount; i++ {
			f1 := ds.Units[l][i].r.Mess.F1
			f2 := ds.Units[l][i].r.Mess.F2
			y1 := float64(ds.Units[l][i].r.Mess.Y1) / 1000
			y2 := float64(ds.Units[l][i].r.Mess.Y2) / 1000

			html += fmt.Sprintf(`
				<tr id="messwerte" data-on:click="window.location.href = 'unit.html?line=%d&uid=%d'">
					<td id="raum">%s</td>
					<td id="com">%s#%02d</td>
					<td id="f1">%d</td>
					<td id="f2">%d</td>
					<td id="y1">%.2f</td>
					<td id="y2">%.2f</td>
					<td id="ts">-</td>
				</tr>
			`, l, i, line.Ort, line.Url, line.Id[i], f1, f2, y1, y2) + "\n"
		}
	}
	html = fmt.Sprintf(`<tbody id="indextable">%s</tbody>`, html)
	return html
}
