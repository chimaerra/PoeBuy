package connections

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"poebuy/modules/connections/headers"
	"poebuy/modules/connections/models"
	"regexp"
)

var ErrorBadPoessid = errors.New("can't get trade info, check POESSID")

const (
	_PcLeagueId = "pc"
)

func GetTradeInfo(poesessid string) (*models.TradeInfo, error) {

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, "https://www.pathofexile.com/trade", nil)
	if err != nil {
		return nil, err
	}
	req.Header = headers.GetFetchitemHeaders(poesessid)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	gr, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, err
	}

	bt, err := io.ReadAll(gr)
	if err != nil {
		return nil, err
	}

	info := &models.TradeInfo{}

	leg := make([]models.League, 0, 10)
	match := regexp.MustCompile(`"leagues": (\[[^\]]*\])`).FindSubmatch(bt)
	if match == nil {
		return nil, ErrorBadPoessid
	}
	json.Unmarshal(match[1], &leg)
	for _, l := range leg {
		if l.Realm == _PcLeagueId {
			info.Leagues = append(info.Leagues, l)
		}
	}

	info.Nickname = string(regexp.MustCompile(`account/view-profile/.*?">(.+?)<`).FindSubmatch(bt)[1])

	return info, nil
}
