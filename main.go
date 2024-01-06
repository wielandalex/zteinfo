package main

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		println("Not enough arguments")
		usage(os.Args[0])
		os.Exit(1)
	}

	ip := os.Args[1]
	password := os.Args[2]

	// Don't verify TLS
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	LD, err := getLD(ip)
	if err != nil {
		log.Fatal(err)
	}

	cookie, err := getCookie(ip, password, LD)
	if err != nil {
		log.Fatal(err)
	}

	info, err := getInfo(ip, cookie)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(info)
}

func usage(exe string) {
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Printf("%s ip_address password\n", exe)
	fmt.Println()
}

func hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}

func getLD(ip string) (string, error) {
	req, err := http.NewRequest("GET", "https://"+ip+"/goform/goform_get_cmd_process?isTest=false&cmd=LD", nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Referer", "https://"+ip+"/")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var LDresp struct {
		LD string `json:"LD"`
	}
	if err := json.Unmarshal(body, &LDresp); err != nil {
		return "", err
	}
	return LDresp.LD, nil
}

func getCookie(ip string, password string, LD string) (*http.Cookie, error) {
	h := strings.ToUpper(hash(password))
	h = strings.ToUpper(hash(h + LD))

	p := url.Values{
		"isTest":   {"false"},
		"goformId": {"LOGIN"},
		"password": {h},
	}
	req, err := http.NewRequest("POST", "https://"+ip+"/goform/goform_set_cmd_process", strings.NewReader(p.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Referer", "https://"+ip+"/")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		if cookie.Name != "stok" {
			continue
		}
		return cookie, nil
	}

	return nil, errors.New("cookie not found")
}

func getInfo(ip string, cookie *http.Cookie) (string, error) {
	req, err := http.NewRequest("GET", "https://"+ip+"/goform/goform_get_cmd_process?isTest=false&multi_data=1&cmd=wa_inner_version,cr_version,network_type,rssi,rscp,rmcc,rmnc,enodeb_id,lte_rsrq,lte_rsrp,Z5g_snr,Z5g_rsrp,ZCELLINFO_band,Z5g_dlEarfcn,lte_ca_pcell_arfcn,lte_ca_pcell_band,lte_ca_scell_band,lte_ca_pcell_bandwidth,lte_ca_scell_info,lte_ca_scell_bandwidth,wan_lte_ca,lte_pci,Z5g_CELL_ID,Z5g_SINR,cell_id,wan_lte_ca,lte_ca_pcell_band,lte_ca_pcell_bandwidth,lte_ca_scell_band,lte_ca_scell_bandwidth,lte_ca_pcell_arfcn,lte_ca_scell_arfcn,lte_multi_ca_scell_info,wan_active_band,nr5g_pci,nr5g_action_band,nr5g_cell_id,lte_snr,ecio,wan_active_channel,nr5g_action_channel,ngbr_cell_info,monthly_tx_bytes,monthly_rx_bytes,lte_pci,lte_pci_lock,lte_earfcn_lock,wan_ipaddr,wan_apn,pm_sensor_mdm,pm_modem_5g,nr5g_pci,nr5g_action_channel,nr5g_action_band,Z5g_SINR,Z5g_rsrp,wan_active_band,wan_active_channel,wan_lte_ca,lte_multi_ca_scell_info,cell_id,dns_mode,prefer_dns_manual,standby_dns_manual,network_type,rmcc,rmnc,lte_rsrq,lte_rssi,lte_rsrp,lte_snr,wan_lte_ca,lte_ca_pcell_band,lte_ca_pcell_bandwidth,lte_ca_scell_band,lte_ca_scell_bandwidth,lte_ca_pcell_arfcn,lte_ca_scell_arfcn,wan_ipaddr,static_wan_ipaddr,opms_wan_mode,opms_wan_auto_mode,ppp_status,loginfo,realtime_time,signalbar,realtime_rx_thrpt,realtime_tx_thrpt", nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(cookie)
	req.Header.Add("Referer", "https://"+ip+"/index.html")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
