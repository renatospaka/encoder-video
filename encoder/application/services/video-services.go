package services

import (
	"encoder/application/repositories"
	"encoder/domain"
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}

	bkt := client.Bucket(bucketName)
	obj := bkt.Object(v.Video.FilePath)
	
	r, err := obj.NewReader(ctx)
	if err != nil {
		return err
	}
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	file, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")
	if err != nil {
		return err
	}

	// test if it worked
	_, err = file.Write(body)
	if err != nil {
		return err
	}

	defer file.Close()

	log.Printf("The video %v has been stored.", v.Video.ID)

	//no errors ocurred
	return nil
}

func (v *VideoService) Fragment() error {
	err := os.Mkdir(os.Getenv("localStoragePath") + "/" + v.Video.ID, os.ModePerm)
	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	printOutput(output)
	return nil
}

func printOutput(out []byte) {
	if len(out) > 0 {
		log.Printf("=====> Output: %s\n", string(out))
	}
}