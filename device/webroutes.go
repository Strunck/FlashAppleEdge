package device

import (
	"fmt"
	"net/http"

	"github.com/starfederation/datastar-go/datastar"
)

func HandleIndexTbl(ds State) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		sse := datastar.NewSSE(w, r)

		err := sse.PatchElements(`
		    <code id="logmessages">
                WUNDERBAR
            </code>
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = sse.PatchElements(loadIndexTbl(ds))
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

			html += fmt.Sprintf(`
				<tr id="messwerte" data-on:click="window.location.href = 'unit.html?line=%d&uid=%d'">
					<td id="raum">%s</td>
					<td id="com">%s#%02d</td>
					<td id="f1">0.00</td>
					<td id="f2">0.00</td>
					<td id="y1">0.00</td>
					<td id="y2">0.00</td>
					<td id="ts">-</td>
				</tr>
			`, l, i, line.Ort, line.Url, line.Id[i]) + "\n"
		}
	}
	html = fmt.Sprintf(`<tbody id="indextable">%s</tbody>`, html)
	return html
}
