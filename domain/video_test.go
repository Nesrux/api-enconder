package domain_test

import (
	"github.com/Nesrux/api-enconder/domain"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Validated_ifVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)

}
