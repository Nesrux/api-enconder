package domain_test

import (
	"testing"
	"time"

	"github.com/Nesrux/api-enconder/domain"

	"github.com/stretchr/testify/require"
)

func Test_Validated_ifVideoIsEmpty(t *testing.T) {
	video := domain.NewVideo()
	err := video.Validate()

	require.Error(t, err)
}

func Test_VideoId_IsNotUUID(t *testing.T) {
	video := domain.NewVideo()
	//when
	video.ID = "abc"
	video.ResourceID = "123"
	video.FilePath = "path"
	video.CreatedAt = time.Now()
	
	//then 
	err := video.Validate()

	require.Error(t, err)
}
