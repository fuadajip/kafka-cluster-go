package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"github.com/fuadajip/kafka-cluster-go/models"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo"
)

// CreateUserUploadPresigned return S3 presigned URL for user's POST
func CreateUserUploadPresigned(c echo.Context, payload *models.CreateUserUploadPresignedRequest) (*models.CreateUserUploadPresignedResponse, error) {
	ac := c.(CustomApplicationContext)
	conf := ac.SharedConf

	sess, err := GetAWSSession(conf)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	currentDate := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	trimedFilename := strings.Replace(*payload.Filename, " ", "", -1)
	slashRemovedFilename := strings.Replace(trimedFilename, "/", "-", -1)
	capitalizedDocumentType := strings.ToUpper(*payload.DocumentType)
	capitalizedUserType := strings.ToUpper(*payload.UserType)

	uniqueFileName := fmt.Sprintf("%s-%s-%d-%s", capitalizedDocumentType, currentDate, t.Unix(), slashRemovedFilename)
	fileKey := fmt.Sprintf("/uploads/private/%s/%s", capitalizedUserType, uniqueFileName)
	fileContentType := mime.TypeByExtension(filepath.Ext(*payload.Filename))

	logger.Info(fileKey)

	svc := s3.New(sess)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(conf.GetAWSBucket()),
		Key:    aws.String(fileKey),
	})

	strURL, err := req.Presign(3 * time.Hour)
	if err != nil {
		return nil, err
	}

	mappedResp := &models.CreateUserUploadPresignedResponse{
		Key:          aws.String(fileKey),
		URL:          aws.String(strURL),
		Filename:     aws.String(uniqueFileName),
		DocumentType: aws.String(capitalizedDocumentType),
		MimeType:     aws.String(fileContentType),
		Size:         payload.Size,
		UserType:     aws.String(capitalizedUserType),
	}

	return mappedResp, nil
}

// CreateUserUploadPresignedView return S3 presigned URL for user's to read
func CreateUserUploadPresignedView(c echo.Context, payload *models.CreateUserUploadPresignedViewRequest) (*models.CreateUserUploadPresignedViewResponse, error) {
	ac := c.(*CustomApplicationContext)
	conf := ac.SharedConf

	sess, err := GetAWSSession(conf)
	if err != nil {
		return nil, err
	}
	svc := s3.New(sess)
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(conf.GetAWSBucket()),
		Key:    payload.Key,
	})

	strURL, err := req.Presign(3 * time.Hour)
	if err != nil {
		return nil, err
	}

	mappedResp := &models.CreateUserUploadPresignedViewResponse{
		URL: aws.String(strURL),
		Key: payload.Key,
	}

	return mappedResp, nil
}

// GetS3ObjectBuffer will get object from s3 and return the object as buffer
func GetS3ObjectBuffer(c echo.Context, key *string) ([]byte, error) {

	ac := c.(*CustomApplicationContext)
	conf := ac.SharedConf

	sess, err := GetAWSSession(conf)
	if err != nil {
		return nil, err
	}

	svc := s3.New(sess)

	req, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(conf.GetAWSBucket()),
		Key:    key,
	})
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	return body, nil
}

// UploadS3Object upload to S3 from base64 string data
func UploadS3Object(c echo.Context, key string, contentType string, body string) error {
	ac := c.(*CustomApplicationContext)
	conf := ac.SharedConf

	sess, err := GetAWSSession(conf)
	if err != nil {
		return err
	}

	data, err := base64.StdEncoding.DecodeString(body)
	if err != nil {
		logger.Error(fmt.Printf("Error uploading %s err: %s", key, err.Error()))
		return err
	}

	svc := s3.New(sess)
	_, err = svc.PutObject(&s3.PutObjectInput{
		ACL:         aws.String("private"),
		Body:        aws.ReadSeekCloser(bytes.NewReader(data)),
		Bucket:      aws.String(conf.GetAWSBucket()),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return err
	}

	return nil
}
