package testsync

import (
	"context"
	"github.com/etcd-carry/etcd-carry/pkg/testing/testcodec"
	"go.etcd.io/etcd/client/v3"
	"testing"
)

func PrepareSequentialTestData(t *testing.T, client *clientv3.Client) {
	if _, err := client.Put(context.TODO(), string(testcodec.SampleNamespaceMatchedKey1), string(testcodec.SampleNamespaceMatchedValue1)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleNamespaceMatchedKey2), string(testcodec.SampleNamespaceMatchedValue2)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleNamespaceMismatchedKey1), string(testcodec.SampleNamespaceMismatchedValue1)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleCrdKey), string(testcodec.SampleCrdValue)); err != nil {
		t.Fatal(err)
	}
}

func PrepareSecondaryTestData(t *testing.T, client *clientv3.Client) {
	if _, err := client.Put(context.TODO(), string(testcodec.SampleConfigmapMatchedKey1), string(testcodec.SampleConfigmapMatchedValue1)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleConfigmapMismatchedKey1), string(testcodec.SampleConfigmapMismatchedValue1)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleSecretMatchedKey1), string(testcodec.SampleSecretMatchedValue1)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleSecretMatchedKey2), string(testcodec.SampleSecretMatchedValue2)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleSecretMismatchedKey1), string(testcodec.SampleSecretMismatchedValue1)); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Put(context.TODO(), string(testcodec.SampleSecretMismatchedKey2), string(testcodec.SampleSecretMismatchedValue2)); err != nil {
		t.Fatal(err)
	}
}

func PrepareTestData(t *testing.T, client *clientv3.Client) {
	PrepareSequentialTestData(t, client)
	PrepareSecondaryTestData(t, client)
}
