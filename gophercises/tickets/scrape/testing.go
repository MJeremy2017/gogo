package scrape

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func setUp() *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data := `
<!doctype html>
<html>
<ul class="prinav">
        <li class="t cat3">
            <a href="/sg/Concert-Tickets">Concert Tickets</a>
        </li>
        <li class="t cat2">
            <a href="/sg/Sports-Tickets">Sports Tickets</a>
        </li>
        <li class="t cat1">
            <a href="/sg/Theatre-Tickets">Theatre Tickets</a>
        </li>
        <li class="t cat1023">
            <a href="/sg/Festival-Tickets">Festival Tickets</a>
        </li>
</ul>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})
	mux.HandleFunc("/sg/Concert-Tickets", func(w http.ResponseWriter, r *http.Request) {
		data := `
<!doctype html>
<html>
<div class="bMag3 pbxxxl pb0-s">
        <div class="pgw">
            <h1 class="mbxs ibk t cWht xxl xl-s ">Concert Tickets</h1>
            <ul class="cloud">
                    <li><a href="/sg/Concert-Tickets/Clubs-and-Dance">Club and dance </a></li>
                    <li><a href="/sg/Theater-Tickets/Flamenco">Flamenco </a></li>
            </ul>
        </div>
</div>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})

	mux.HandleFunc("/sg/Concert-Tickets/Clubs-and-Dance", func(w http.ResponseWriter, r *http.Request) {
		data := `
<html>
    <div class="uuxxl pgw">
		<a href="something">abc</a>
        <h2 class="txtc t xxl pbl">All Events</h2>
            <h3 class="t xxl"></h3>
            <ul class="cloud mbxl">
                    <li>
                        <a class="t" href="/sg/Concert-Tickets/Rock-and-Pop/5-Seconds-of-Summer-Tickets">5 Seconds of Summer</a>
                    </li>
            </ul>
            <h3 class="t xxl">B</h3>
            <ul class="cloud mbxl">
                    <li>
                        <a class="t" href="/sg/Concert-Tickets/Rock-and-Pop/Bastille-Tickets">Bastille</a>
                    </li>
                    <li>
                        <a class="t" href="/sg/Concert-Tickets/Rock-and-Pop/Brigitte-Tickets">Brigitte</a>
                    </li>
            </ul>
    </div>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})

	mux.HandleFunc("/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets", func(w http.ResponseWriter, r *http.Request) {
		data := `
<html>
                <div class="tablet-header-category-title h l prxxs">
                            <span>
                                <a href="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets" id="catNameInHeader" class="h xxl rel cWht">
                                    Super Junior Tickets
                                </a>
                            </span>

                </div>
<a href="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327" target="_blank" class="js-event-row-container el-row-anchor cGry1">
    <div data-new-tab-nav="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327" data-id="151336327" class="el-row-div js-event-item bdrt">
        <div class="t s ins bdrr date-div pt0 pbm">
            <div>
                <div class="count-down-sash uuxxs">In <span class="fs14">17</span> Days</div>
                <time datetime="2023-02-09T20:00:00">
                    <div class="mbxxs-m">
                        <span class="h m s-l">Thursday, February 09</span>
                    </div>
                    <div class="cGry2 fs14">
                        <i class="i-time vmid"></i>
                        <span class="vmid">20:00</span>
                    </div>
                </time>
            </div>
        </div>

        <div class="el-column-info uum ins">
            <div>
                <div class="t s">
                    <i class="i-location vmid fs14"></i>
                    Espaço Unimed (formerly Espaco das Americas), <span class="t-b">Sao Paulo, Brazil</span>
                </div>
                <span class="t s">
                    <span class="camo bk" href="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327" title="Super Junior at Espa&#231;o Unimed (formerly Espaco das Americas) Sao Paulo on Thu 09 Feb 2023 20:00">
                        <strong class="cGry3">Super Junior</strong>
                        <span class="label-vg-deal bBlu1 ins fs12 txtc cWht"> <span class="i-eye rel dealbadge-star mrxxs"></span>Recently Viewed </span>
                    </span>
                </span>
            </div>
        </div>

    </div>
</a>
<a href="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151396140" target="_blank" class="js-event-row-container el-row-anchor cGry1">
    <div data-new-tab-nav="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151396140" data-id="151396140" class="el-row-div js-event-item">
        <div class="t s ins bdrr date-div pt0 pbm">
            <div>
                <div class="count-down-sash uuxxs">In <span class="fs14">25</span> Days</div>
                <time datetime="2023-02-18T19:30:00">
                    <div class="mbxxs-m">
                        <span class="h m s-l">Saturday, February 18</span>
                    </div>
                    <div class="cGry2 fs14">
                        <i class="i-time vmid"></i>
                        <span class="vmid">19:30</span>
                    </div>
                </time>
				<time class="hide" datetime="2023-02-20T23:59:00"></time>
            </div>
        </div>

        <div class="el-column-info uum ins">
            <div>
                <div class="t s">
                    <i class="i-location vmid fs14"></i>
                    Sunway Lagoon, <span class="t-b">Selangor, Malaysia</span>
                </div>
                <span class="t s">
                    <span class="camo bk" href="/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151396140" title="Kpop Mega Concert at Sunway Lagoon Selangor on Sat 18 Feb 2023 19:30">
                        <strong class="cGry3">Kpop Mega Concert</strong>
                    </span>
                </span>

                <div class="h txtl flex mtxxs flex-middle">
                    <i class="i-stadium vmid mrxxs fs16 nudge-t-1"></i>
                    <span class="vmid">Venue capacity: 2000</span>
                </div>
            </div>
        </div>
	</div>
</a>
</html>
`
		_, _ = fmt.Fprint(w, data)
	})

	mux.HandleFunc("/sg/Concert-Tickets/Rock-and-Pop/Super-Junior-Tickets/E-151336327", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			_, _ = fmt.Fprint(w, "")
		}

		data := `
{
	"Items": [
		{
			"Id": 5843255783,
      		"EventId": 151336327,
			"RawPrice": 108.82,
 			"TicketsLeftInListingMessage": {
				"Message": "4 tickets remaining",
				"Qualifier": "in this listing on our site",
				"Disclaimer": null,
				"HasValue": true,
				"FeatureTrackingKey": "2fbb0992-0355-42c8-98d2-4f0616842abb"
		  	},
			"QuantityRange": "1 - 4",
			"BuyUrl": "/buy/catId=123"
		},
		{
			"Id": 5891820447,
      		"EventId": 151336327,
			"RawPrice": 206.64,
			"TicketsLeftInListingMessage": null,
			"QuantityRange": "1",
			"BuyUrl": "/buy/catId=124"
		}
	]
}
`
		_, _ = fmt.Fprint(w, data)
	})

	return httptest.NewServer(mux)
}
