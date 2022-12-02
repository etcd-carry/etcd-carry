package schema

import (
	"bytes"
	"go.etcd.io/etcd/mvcc/mvccpb"
	mirrorcontext "orcastack.io/etcd-mirror/pkg/mirror/context"
	"orcastack.io/etcd-mirror/pkg/testing/testcodec"
	"orcastack.io/etcd-mirror/pkg/testing/testoptions"
	"testing"
)

func TestAuthenticatedDataString_AuthenticatedData(t *testing.T) {
	testKey := authenticatedDataString(testcodec.SampleNamespaceMatchedKey1)
	if !bytes.Equal(testKey.AuthenticatedData(), testcodec.SampleNamespaceMatchedKey1) {
		t.Fatal("should not happened")
	}
}

func TestDecode(t *testing.T) {
	o := testoptions.GetMirrorOptions(t)
	ctx, err := mirrorcontext.NewMirrorContext(o)
	if err != nil {
		t.Fatal(err)
	}
	testCase := []struct {
		name         string
		key          []byte
		value        []byte
		expectErrors bool
		isMatched    bool
	}{
		{
			name:         "valid and matched Namespace object1",
			key:          testcodec.SampleNamespaceMatchedKey1,
			value:        testcodec.SampleNamespaceMatchedValue1,
			expectErrors: false,
			isMatched:    true,
		},
		{
			name:         "valid but mismatched Namespace object1",
			key:          testcodec.SampleNamespaceMismatchedKey1,
			value:        testcodec.SampleNamespaceMismatchedValue1,
			expectErrors: false,
			isMatched:    false,
		},
		{
			name:         "valid and matched Configmap object1",
			key:          testcodec.SampleConfigmapMatchedKey1,
			value:        testcodec.SampleConfigmapMatchedValue1,
			expectErrors: false,
			isMatched:    true,
		},
		{
			name:         "valid and mismatched Configmap object1",
			key:          testcodec.SampleConfigmapMismatchedKey1,
			value:        testcodec.SampleConfigmapMismatchedValue1,
			expectErrors: false,
			isMatched:    false,
		},
		{
			name:         "valid and matched Secret object1",
			key:          testcodec.SampleSecretMatchedKey1,
			value:        testcodec.SampleSecretMatchedValue1,
			expectErrors: false,
			isMatched:    true,
		},
		{
			name:         "valid and mismatched Secret object1",
			key:          testcodec.SampleSecretMismatchedKey1,
			value:        testcodec.SampleSecretMismatchedValue1,
			expectErrors: false,
			isMatched:    false,
		},
		{
			name:         "valid CRD object",
			key:          testcodec.SampleCrdKey,
			value:        testcodec.SampleCrdValue,
			expectErrors: false,
			isMatched:    true,
		},
		{
			name:         "valid CRD resource object",
			key:          testcodec.SampleCrdResourceKey,
			value:        testcodec.SampleCrdResourceValue,
			expectErrors: false,
			isMatched:    false,
		},
		{
			name:         "valid Service resource object1",
			key:          testcodec.SampleServiceMatchedKey1,
			value:        testcodec.SampleServiceMatchedValue1,
			expectErrors: false,
			isMatched:    true,
		},
		{
			name:         "invalid k8s object1",
			key:          []byte("test1"),
			value:        []byte("{TestString: \"foo\"}"),
			expectErrors: true,
			isMatched:    false,
		},
		{
			name:         "invalid k8s object2",
			key:          []byte("test2"),
			value:        []byte("{\"TestString\": \"foo\"}"),
			expectErrors: true,
			isMatched:    false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			rs, err := Decode(ctx, tc.key, tc.value)
			if err != nil && !tc.expectErrors {
				t.Errorf("expect no errors, but error found %+v", err)
			}
			if err == nil && tc.expectErrors {
				t.Errorf("expect errors, but no error found %+v", err)
			}
			if err == nil && !tc.expectErrors {
				if (rs.FilterSequentialByRules(ctx) || rs.FilterSecondaryByRules(ctx)) != tc.isMatched {
					t.Errorf("should match as expected, but does not")
				}

				_, err = rs.MutateKubeResource(&mvccpb.KeyValue{
					Key:   tc.key,
					Value: tc.value,
				})
				if err != nil {
					t.Errorf("should not happened, found errors %+v", err)
				}
			}
		})
	}
}
