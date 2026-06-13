package serverV1

import (
	"fmt"
	"net/http"

	"github.com/mohamedveron/go_app_template/proxy"
)

func (*HTTP) GetParagraphByTopic(w http.ResponseWriter, r *http.Request, topic string) {
	gpt := proxy.NewOpenAI("token")
	msg := gpt.GetMessage(topic)
	fmt.Print(msg)
	w.Write([]byte(msg))
}
