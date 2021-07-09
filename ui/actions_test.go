package ui

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tenntenn/natureremo"
)

func NewMockServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Error(err)
			return
		}
		handler(w, r)
	}))

	return server
}

func NewMockClient(t *testing.T, url string) *natureremo.Client {
	c := natureremo.NewClient("tmp")
	c.BaseURL = url
	return c
}

func TestActionAppliancePower(t *testing.T) {
	state := &State{
		Appliances: []*natureremo.Appliance{
			{
				ID: "1",
				Light: &natureremo.Light{
					State: &natureremo.LightState{
						Power: "off",
					},
				},
				TV: &natureremo.TV{
					State: &natureremo.TVState{
						Input: "",
					},
				},
				AirConSettings: &natureremo.AirConSettings{
					Temperature:   "25",
					OperationMode: "cool",
					AirVolume:     "4",
					AirDirection:  "swing",
					Button:        "",
				},
			},
		},
	}

	t.Run("power on/off", func(t *testing.T) {
		tests := []struct {
			appType natureremo.ApplianceType
			power   natureremo.Button
			want    string
		}{
			{
				appType: natureremo.ApplianceTypeAirCon,
				power:   natureremo.ButtonPowerOn,
				want:    natureremo.ButtonPowerOn.StringValue(),
			},
			{
				appType: natureremo.ApplianceTypeAirCon,
				power:   natureremo.ButtonPowerOff,
				want:    natureremo.ButtonPowerOff.StringValue(),
			},
			{
				appType: natureremo.ApplianceTypeTV,
				power:   natureremo.ButtonPowerOn,
				want:    "power",
			},
			{
				appType: natureremo.ApplianceTypeTV,
				power:   natureremo.ButtonPowerOff,
				want:    "power",
			},
			{
				appType: natureremo.ApplianceTypeLight,
				power:   natureremo.ButtonPowerOn,
				want:    "on",
			},
			{
				appType: natureremo.ApplianceTypeLight,
				power:   natureremo.ButtonPowerOff,
				want:    "off",
			},
		}

		for _, tt := range tests {
			h := func(w http.ResponseWriter, r *http.Request) {
				got := r.Form.Get("button")
				if got != tt.want {
					t.Errorf("unexpected button %s state, want: %v, got: %v", tt.appType, tt.want, got)
				}

				_, _ = w.Write([]byte("null"))
			}
			server := NewMockServer(t, h)
			cli := NewMockClient(t, server.URL)

			state.Appliances[0].Type = tt.appType
			ctx := AppliancePowerOnOff{
				Row:   0,
				Power: tt.power,
			}

			if err := ActionAppliancesPower(state, cli, ctx); err != nil {
				_ = server.Config.Shutdown(context.Background())
				t.Error(err)
			}

			_ = server.Config.Shutdown(context.Background())
		}

	})

	t.Run("invalid ctx", func(t *testing.T) {
		want := `ctx is not "AppliancePowerOnOff": int`
		cli := NewMockClient(t, "")
		ctx := 1

		got := ActionAppliancesPower(state, cli, ctx).Error()
		if want != got {
			t.Errorf("unexpected error, want: %v, got: %v", want, got)
		}
	})

	t.Run("unsupported appliance type", func(t *testing.T) {
		want := `unsupported appliance: x`
		cli := NewMockClient(t, "")
		state.Appliances[0].Type = natureremo.ApplianceType("x")
		ctx := AppliancePowerOnOff{
			Row:   0,
			Power: "",
		}

		got := ActionAppliancesPower(state, cli, ctx).Error()
		if want != got {
			t.Errorf("unexpected error, want: %v, got: %v", want, got)
		}
	})
}

func TestActionUpdateAirConSettings(t *testing.T) {
	tests := []struct {
		ctx   UpdateAirConFormData
		state *State
		want  *State
	}{
		{
			ctx: UpdateAirConFormData{
				Power: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"ON",
					},
				},
				Mode: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"below",
					},
				},
				Volume: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"1",
					},
				},
			},
			state: &State{
				Appliances: []*natureremo.Appliance{
					{
						ID: "1",
						AirConSettings: &natureremo.AirConSettings{
							OperationMode: "below",
							AirVolume:     "0",
						},
					},
				},
			},
			want: &State{
				Appliances: []*natureremo.Appliance{
					{
						ID: "1",
						AirConSettings: &natureremo.AirConSettings{
							OperationMode: "below",
							AirVolume:     "1",
						},
					},
				},
			},
		},
		{
			ctx: UpdateAirConFormData{
				Power: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"OFF",
					},
				},
				Mode: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"cool",
					},
				},
				Volume: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"1",
					},
				},
				Temp: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"10",
					},
				},
			},
			state: &State{
				Appliances: []*natureremo.Appliance{
					{
						ID: "1",
						AirConSettings: &natureremo.AirConSettings{
							OperationMode: "below",
							AirVolume:     "0",
							Temperature:   "2",
						},
					},
				},
			},
			want: &State{
				Appliances: []*natureremo.Appliance{
					{
						ID: "1",
						AirConSettings: &natureremo.AirConSettings{
							OperationMode: "cool",
							AirVolume:     "1",
							Temperature:   "10",
							Button:        natureremo.ButtonPowerOff,
						},
					},
				},
			},
		},
		{
			ctx: UpdateAirConFormData{
				Power: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"OFF",
					},
				},
				Mode: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"dry",
					},
				},
				Temp: UpdateAirConFormItem{
					Current: 0,
					Values: []string{
						"10",
					},
				},
			},
			state: &State{
				Appliances: []*natureremo.Appliance{
					{
						ID: "1",
						AirConSettings: &natureremo.AirConSettings{
							OperationMode: "below",
							Temperature:   "2",
						},
					},
				},
			},
			want: &State{
				Appliances: []*natureremo.Appliance{
					{
						ID: "1",
						AirConSettings: &natureremo.AirConSettings{
							OperationMode: "dry",
							Temperature:   "10",
							Button:        natureremo.ButtonPowerOff,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		h := func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte("null"))
		}
		server := NewMockServer(t, h)
		cli := NewMockClient(t, server.URL)

		tt.state.SelectApplianceIdx = 0
		if err := ActionUpdateAirConSettings(tt.state, cli, tt.ctx); err != nil {
			_ = server.Config.Shutdown(context.Background())
			t.Fatal(err)
		}

		if diff := cmp.Diff(tt.want, tt.state); diff != "" {
			_ = server.Config.Shutdown(context.Background())
			t.Fatalf("the state has diff (-want +got):\n%v", diff)
		}
		_ = server.Config.Shutdown(context.Background())
	}

}
