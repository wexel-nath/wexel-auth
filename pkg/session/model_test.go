package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewSessionFromRow(t *testing.T) {
	now := time.Now().Unix()

	type want struct{
		session Session
		wantErr bool
	}
	tests := map[string]struct{
		row  map[string]interface{}
		want
	}{
		"success": {
			row:  map[string]interface{}{
				columnSessionID: "fake.session.id",
				columnUserID:    int64(1),
				columnTimestamp: now,
				columnExpiry:    now + 300,
			},
			want: want{
				session: Session{
					SessionID: "fake.session.id",
					UserID:    1,
					Timestamp: now,
					Expiry:    now + 300,
				},
				wantErr: false,
			},
		},
		"missing session id": {
			row:  map[string]interface{}{
				columnUserID:    int64(1),
				columnTimestamp: now,
				columnExpiry:    now + 300,
			},
			want: want{
				wantErr: true,
			},
		},
		"missing user id": {
			row:  map[string]interface{}{
				columnSessionID: "fake.session.id",
				columnTimestamp: now,
				columnExpiry:    now + 300,
			},
			want: want{
				wantErr: true,
			},
		},
		"missing timestamp": {
			row:  map[string]interface{}{
				columnSessionID: "fake.session.id",
				columnUserID:    int64(1),
				columnExpiry:    now + 300,
			},
			want: want{
				wantErr: true,
			},
		},
		"missing expiry": {
			row:  map[string]interface{}{
				columnSessionID: "fake.session.id",
				columnUserID:    int64(1),
				columnTimestamp: now,
			},
			want: want{
				wantErr: true,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(st *testing.T) {
			session, err := newSessionFromRow(test.row)

			assert.Equal(t, test.wantErr, err != nil)
			if !test.wantErr {
				assert.Equal(t, test.session, session)
			}
		})
	}
}
