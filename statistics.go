package bigip

import (
	"encoding/json"
	"fmt"
	"strings"
)

//Disgusting struct for pool statistics
//{
//  "entries": {
//    "activeMemberCnt": {
//      "value": 0
//    },
type PoolStatistics struct {
	Entries struct {
		Activemembercnt struct {
			Value float64 `json:"value"`
		} `json:"activeMemberCnt"`
		ConnqAgeedm struct {
			Value float64 `json:"value"`
		} `json:"connq.ageEdm"`
		ConnqAgeema struct {
			Value float64 `json:"value"`
		} `json:"connq.ageEma"`
		ConnqAgehead struct {
			Value float64 `json:"value"`
		} `json:"connq.ageHead"`
		ConnqAgemax struct {
			Value float64 `json:"value"`
		} `json:"connq.ageMax"`
		ConnqDepth struct {
			Value float64 `json:"value"`
		} `json:"connq.depth"`
		ConnqServiced struct {
			Value float64 `json:"value"`
		} `json:"connq.serviced"`
		ConnqallAgeedm struct {
			Value float64 `json:"value"`
		} `json:"connqAll.ageEdm"`
		ConnqallAgeema struct {
			Value float64 `json:"value"`
		} `json:"connqAll.ageEma"`
		ConnqallAgehead struct {
			Value float64 `json:"value"`
		} `json:"connqAll.ageHead"`
		ConnqallAgemax struct {
			Value float64 `json:"value"`
		} `json:"connqAll.ageMax"`
		ConnqallDepth struct {
			Value float64 `json:"value"`
		} `json:"connqAll.depth"`
		ConnqallServiced struct {
			Value float64 `json:"value"`
		} `json:"connqAll.serviced"`
		Cursessions struct {
			Value float64 `json:"value"`
		} `json:"curSessions"`
		Minactivemembers struct {
			Value float64 `json:"value"`
		} `json:"minActiveMembers"`
		Monitorrule struct {
			Description string `json:"description"`
		} `json:"monitorRule"`
		ServersideBitsin struct {
			Value float64 `json:"value"`
		} `json:"serverside.bitsIn"`
		ServersideBitsout struct {
			Value float64 `json:"value"`
		} `json:"serverside.bitsOut"`
		ServersideCurconns struct {
			Value float64 `json:"value"`
		} `json:"serverside.curConns"`
		ServersideMaxconns struct {
			Value float64 `json:"value"`
		} `json:"serverside.maxConns"`
		ServersidePktsin struct {
			Value float64 `json:"value"`
		} `json:"serverside.pktsIn"`
		ServersidePktsout struct {
			Value float64 `json:"value"`
		} `json:"serverside.pktsOut"`
		ServersideTotconns struct {
			Value float64 `json:"value"`
		} `json:"serverside.totConns"`
		StatusAvailabilitystate struct {
			Description string `json:"description"`
		} `json:"status.availabilityState"`
		StatusEnabledstate struct {
			Description string `json:"description"`
		} `json:"status.enabledState"`
		StatusStatusreason struct {
			Description string `json:"description"`
		} `json:"status.statusReason"`
		Tmname struct {
			Description string `json:"description"`
		} `json:"tmName"`
		Totrequests struct {
			Value float64 `json:"value"`
		} `json:"totRequests"`
	} `json:"entries"`
	Generation int64    `json:"generation"`
	Kind       string `json:"kind"`
	Selflink   string `json:"selfLink"`
}

// Distgusting struct for VSS stats
//{
//  "entries": {
//    "clientside.bitsIn": {
//      "value": 4948143240
//    },
type VSSStatistics struct {
	Entries struct {
		ClientsideBitsin struct {
			Value float64 `json:"value"`
		} `json:"clientside.bitsIn"`
		ClientsideBitsout struct {
			Value float64 `json:"value"`
		} `json:"clientside.bitsOut"`
		ClientsideCurconns struct {
			Value float64 `json:"value"`
		} `json:"clientside.curConns"`
		ClientsideMaxconns struct {
			Value float64 `json:"value"`
		} `json:"clientside.maxConns"`
		ClientsidePktsin struct {
			Value float64 `json:"value"`
		} `json:"clientside.pktsIn"`
		ClientsidePktsout struct {
			Value float64 `json:"value"`
		} `json:"clientside.pktsOut"`
		ClientsideTotconns struct {
			Value float64 `json:"value"`
		} `json:"clientside.totConns"`
		Cmpenablemode struct {
			Description string `json:"description"`
		} `json:"cmpEnableMode"`
		Cmpenabled struct {
			Description string `json:"description"`
		} `json:"cmpEnabled"`
		Csmaxconndur struct {
			Value float64 `json:"value"`
		} `json:"csMaxConnDur"`
		Csmeanconndur struct {
			Value float64 `json:"value"`
		} `json:"csMeanConnDur"`
		Csminconndur struct {
			Value float64 `json:"value"`
		} `json:"csMinConnDur"`
		Destination struct {
			Description string `json:"description"`
		} `json:"destination"`
		EphemeralBitsin struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.bitsIn"`
		EphemeralBitsout struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.bitsOut"`
		EphemeralCurconns struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.curConns"`
		EphemeralMaxconns struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.maxConns"`
		EphemeralPktsin struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.pktsIn"`
		EphemeralPktsout struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.pktsOut"`
		EphemeralTotconns struct {
			Value float64 `json:"value"`
		} `json:"ephemeral.totConns"`
		Fiveminavgusageratio struct {
			Value float64 `json:"value"`
		} `json:"fiveMinAvgUsageRatio"`
		Fivesecavgusageratio struct {
			Value float64 `json:"value"`
		} `json:"fiveSecAvgUsageRatio"`
		Oneminavgusageratio struct {
			Value float64 `json:"value"`
		} `json:"oneMinAvgUsageRatio"`
		StatusAvailabilitystate struct {
			Description string `json:"description"`
		} `json:"status.availabilityState"`
		StatusEnabledstate struct {
			Description string `json:"description"`
		} `json:"status.enabledState"`
		StatusStatusreason struct {
			Description string `json:"description"`
		} `json:"status.statusReason"`
		SyncookieAccepts struct {
			Value float64 `json:"value"`
		} `json:"syncookie.accepts"`
		SyncookieHwaccepts struct {
			Value float64 `json:"value"`
		} `json:"syncookie.hwAccepts"`
		SyncookieHwsyncookies struct {
			Value float64 `json:"value"`
		} `json:"syncookie.hwSyncookies"`
		SyncookieHwsyncookieinstance struct {
			Value float64 `json:"value"`
		} `json:"syncookie.hwsyncookieInstance"`
		SyncookieRejects struct {
			Value float64 `json:"value"`
		} `json:"syncookie.rejects"`
		SyncookieSwsyncookieinstance struct {
			Value float64 `json:"value"`
		} `json:"syncookie.swsyncookieInstance"`
		SyncookieSyncachecurr struct {
			Value float64 `json:"value"`
		} `json:"syncookie.syncacheCurr"`
		SyncookieSyncacheover struct {
			Value float64 `json:"value"`
		} `json:"syncookie.syncacheOver"`
		SyncookieSyncookies struct {
			Value float64 `json:"value"`
		} `json:"syncookie.syncookies"`
		Syncookiestatus struct {
			Description string `json:"description"`
		} `json:"syncookieStatus"`
		Tmname struct {
			Description string `json:"description"`
		} `json:"tmName"`
		Totrequests struct {
			Value float64 `json:"value"`
		} `json:"totRequests"`
	} `json:"entries"`
	Generation int64    `json:"generation"`
	Kind       string `json:"kind"`
	Selflink   string `json:"selfLink"`
}

// Get pool statistics by full path(/<partition>/<name>)
func (b *BigIP) GetPoolStatistics(fullpath string) (*PoolStatistics, error) {
	fixedpath := strings.Replace(fullpath, "/", "~", -1)
	resp, err := b.SafeGet(fmt.Sprintf("%s/%s/stats", uriPool, fixedpath))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var stats PoolStatistics
	err = json.Unmarshal(resp, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// Get pool statistics by full path (/<partition>/<name>)
func (b *BigIP) GetVSSStatistics(fullpath string) (*VSSStatistics, error) {
	fixedpath := strings.Replace(fullpath, "/", "~", -1)
	resp, err := b.SafeGet(fmt.Sprintf("%s/%s/stats", uriVirtual, fixedpath))
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, nil
	}

	var stats VSSStatistics
	err = json.Unmarshal(resp, &stats)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
