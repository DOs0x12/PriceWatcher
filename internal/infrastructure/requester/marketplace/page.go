package marketplace

import (
	"PriceWatcher/internal/entities/page"
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

const headlessCom = `google-chrome --headless=new --dump-dom --no-sandbox --window-size=1920x1080--user-agent='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36' --lang=ru-RU --run-all-compositor-stages-before-draw --virtual-time-budget=5000 --timeout=5000 --use-gl --disable-gpu`

type MarketplaceRequester struct{}

func (r MarketplaceRequester) RequestPage(url string) (page.Response, error) {
	cmd := exec.Command("sh", "-c", headlessCom+" '"+url+"'")
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return page.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", err)
	}
	if outb.Len() == 0 {
		return page.Response{Body: nil}, fmt.Errorf("cannot get the data from the address: %v", errb.String())
	}
	reader := strings.NewReader(outb.String())

	return page.Response{Body: reader}, nil
}
