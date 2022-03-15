package request

import (
	"github.com/coredns/coredns/plugin/test"
	"github.com/miekg/dns"
	"reflect"
	"testing"
)

func TestNewScrubWriter(t *testing.T) {
	type args struct {
		req *dns.Msg
		w   dns.ResponseWriter
	}
	tests := []struct {
		name string
		args args
		want *ScrubWriter
	}{
		{
			name: "Test NewScrubWriter",
			args: args{
				req: &dns.Msg{
					Question: []dns.Question{
						{
							Name:   "google.com.",
							Qtype:  dns.TypeA,
							Qclass: dns.ClassINET,
						},
					},
				},
				w: nil,
			},
			want: &ScrubWriter{
				req: &dns.Msg{
					Question: []dns.Question{
						{
							Name:   "google.com.",
							Qtype:  dns.TypeA,
							Qclass: dns.ClassINET,
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewScrubWriter(tt.args.req, tt.args.w); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewScrubWriter() = %v, want %v", got, tt.want)
			}
		})
	}
}

type inmemoryWriter struct {
	test.ResponseWriter
	written []byte
}

func (r *inmemoryWriter) WriteMsg(m *dns.Msg) error {
	r.written, _ = m.Pack()
	return r.ResponseWriter.WriteMsg(m)
}

func TestScrubWriter_WriteMsg(t *testing.T) {
	type fields struct {
		ResponseWriter dns.ResponseWriter
		req            *dns.Msg
	}
	type args struct {
		m *dns.Msg
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Test ScrubWriter WriteMsg",
			fields: fields{
				ResponseWriter: &inmemoryWriter{
					ResponseWriter: test.ResponseWriter{},
				},
				req: &dns.Msg{
					Question: []dns.Question{
						{
							Name:   "google.com.",
							Qtype:  dns.TypeA,
							Qclass: dns.ClassINET,
						},
					},
				},
			},
			args: args{
				m: &dns.Msg{
					Question: []dns.Question{
						{
							Name:   "google.com.",
							Qtype:  dns.TypeA,
							Qclass: dns.ClassINET,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScrubWriter{
				ResponseWriter: tt.fields.ResponseWriter,
				req:            tt.fields.req,
			}
			if err := s.WriteMsg(tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("WriteMsg() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
