package db

import (
	"strings"
	"testing"
)

func compareDSN(t *testing.T, got, want string) {
	gotParts := strings.SplitN(got, "?", 2)
	wantParts := strings.SplitN(want, "?", 2)
	if len(gotParts) == 2 && len(wantParts) == 2 {
		if gotParts[0] != wantParts[0] {
			t.Errorf("dsn base = %q, want %q", gotParts[0], wantParts[0])
		}
		gotParams := strings.Split(gotParts[1], "&")
		wantParams := strings.Split(wantParts[1], "&")
		gotMap := map[string]string{}
		wantMap := map[string]string{}
		for _, p := range gotParams {
			kv := strings.SplitN(p, "=", 2)
			if len(kv) == 2 {
				gotMap[kv[0]] = kv[1]
			}
		}
		for _, p := range wantParams {
			kv := strings.SplitN(p, "=", 2)
			if len(kv) == 2 {
				wantMap[kv[0]] = kv[1]
			}
		}
		if len(gotMap) != len(wantMap) {
			t.Errorf("dsn param count = %d, want %d", len(gotMap), len(wantMap))
		}
		for k, v := range wantMap {
			if gotMap[k] != v {
				t.Errorf("dsn param %q = %q, want %q", k, gotMap[k], v)
			}
		}
	} else {
		if got != want {
			t.Errorf("dsn() = %q, want %q", got, want)
		}
	}
}

func TestDBConfig_dsn(t *testing.T) {
	for name, tt := range dsnTestCases {
		t.Run(name, func(t *testing.T) {
			// Arrange
			cfg := tt.config
			want := tt.wantDSN

			// Act
			got := cfg.dsn()

			// Assert
			if strings.HasPrefix(want, "root:@tcp(") || strings.HasPrefix(want, "user:pass@tcp(") {
				compareDSN(t, got, want)
			} else {
				if got != want {
					t.Errorf("dsn() = %q, want %q", got, want)
				}
			}
		})
	}
}

func TestConnect(t *testing.T) {
	for name, tt := range connectTestCases {
		t.Run(name, func(t *testing.T) {
			// Act
			db, err := Connect(tt.config)

			// Assert
			if tt.wantErr {
				if err == nil {
					t.Fatal("Connect() error = nil, wantErr = true")
				}
				if err.Error() != tt.errMsg {
					t.Errorf("Connect() error = %q, want %q", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("Connect() unexpected error = %v", err)
			}
			if db == nil {
				t.Fatal("Connect() returned nil *gorm.DB")
			}
		})
	}
}
