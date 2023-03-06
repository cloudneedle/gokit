package config

import "testing"

func TestEtcd_Put(t *testing.T) {
	cli, err := New(WithType(ETCD), WithPath("127.0.0.1:2379"))
	if err != nil {
		t.Fatal(err)
	}
	defer cli.Close()

	err = cli.Etcd.Put("admin/server_host", ":8081")
	if err != nil {
		t.Fatal(err)
	}
	//err = cli.Etcd.Put("admin/log_path", "app.log")
	//if err != nil {
	//	t.Fatal(err)
	//}
}
