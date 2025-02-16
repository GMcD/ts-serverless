package handlers

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"

	models "github.com/GMcD/ts-serverless/micros/posts/models"
	"github.com/alexellis/hmac"
	"github.com/gofiber/fiber/v2"
	uuid "github.com/gofrs/uuid"
	coreConfig "github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/pkg/log"
	"github.com/red-gold/telar-core/types"
	"github.com/red-gold/telar-core/utils"
)

const contentMaxLength = 20

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// StringRand
func StringRand(length int) string {
	return StringWithCharset(length, charset)
}

// generatePostURLKey
func generatePostURLKey(socialName, body, postId string) string {
	content := body
	if contentMaxLength <= len(body) {
		content = body[:contentMaxLength]
	}

	return strings.ToLower(fmt.Sprintf("%s_%s-post-%s-%s", socialName, strings.ReplaceAll(content, " ", "-"), strings.Split(postId, "-")[0], StringRand(5)))
}

// Async Call Stuff To Propagate Context
type ResultAsync struct {
	Result []byte
	Error  error
}

type UserInfoInReq struct {
	UserId      uuid.UUID `json:"uid"`
	Username    string    `json:"email"`
	DisplayName string    `json:"displayName"`
	SocialName  string    `json:"socialName"`
	Avatar      string    `json:"avatar"`
	Banner      string    `json:"banner"`
	TagLine     string    `json:"tagLine"`
	SystemRole  string    `json:"role"`
	CreatedDate int64     `json:"createdDate"`
}

// getHeadersFromUserInfoReq
func getHeadersFromUserInfoReq(info *UserInfoInReq) map[string][]string {
	userHeaders := make(map[string][]string)
	userHeaders["uid"] = []string{info.UserId.String()}
	userHeaders["email"] = []string{info.Username}
	userHeaders["avatar"] = []string{info.Avatar}
	userHeaders["banner"] = []string{info.Banner}
	userHeaders["tagLine"] = []string{info.TagLine}
	userHeaders["displayName"] = []string{info.DisplayName}
	userHeaders["socialName"] = []string{info.SocialName}
	userHeaders["role"] = []string{info.SystemRole}

	return userHeaders
}

// getUserInfoReq
func getUserInfoReq(c *fiber.Ctx) *UserInfoInReq {
	currentUser, ok := c.Locals("user").(types.UserContext)
	if !ok {
		return &UserInfoInReq{}
	}
	userInfoInReq := &UserInfoInReq{
		UserId:      currentUser.UserID,
		Username:    currentUser.Username,
		Avatar:      currentUser.Avatar,
		DisplayName: currentUser.DisplayName,
		SystemRole:  currentUser.SystemRole,
	}
	return userInfoInReq

}

// getHeaderInfoReq
func getHeaderInfoReq(c *fiber.Ctx) map[string][]string {
	return getHeadersFromUserInfoReq(getUserInfoReq(c))
}

// Async Call Stuff To Propagate Context

// functionCall send request to another function/microservice using cookie validation
func functionCall(method string, bytesReq []byte, url string, header map[string][]string) ([]byte, error) {
	prettyURL := utils.GetPrettyURLf(url)
	bodyReader := bytes.NewBuffer(bytesReq)

	httpReq, httpErr := http.NewRequest(method, *coreConfig.AppConfig.InternalGateway+prettyURL, bodyReader)
	if httpErr != nil {
		return nil, httpErr
	}
	payloadSecret := *coreConfig.AppConfig.PayloadSecret

	digest := hmac.Sign(bytesReq, []byte(payloadSecret))
	httpReq.Header.Set("Content-type", "application/json")
	fmt.Printf("\ndigest: %s, header: %v \n", "sha1="+hex.EncodeToString(digest), types.HeaderHMACAuthenticate)
	httpReq.Header.Add(types.HeaderHMACAuthenticate, "sha1="+hex.EncodeToString(digest))

	if header != nil {
		for k, v := range header {
			httpReq.Header[k] = v
		}
	}

	utils.AddPolicies(httpReq)

	c := http.Client{}
	res, reqErr := c.Do(httpReq)
	fmt.Printf("\nUrl : %s, Result : %v\n", url, *res)
	if reqErr != nil {
		return nil, fmt.Errorf("Error while sending admin check request!: %s", reqErr.Error())
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	resData, readErr := ioutil.ReadAll(res.Body)
	if resData == nil || readErr != nil {
		return nil, fmt.Errorf("failed to read response from admin check request.")
	}

	if res.StatusCode != http.StatusAccepted && res.StatusCode != http.StatusOK {
		if res.StatusCode == http.StatusNotFound {
			return nil, NotFoundHTTPStatusError
		}
		return nil, fmt.Errorf("failed to call %s api, invalid status: %s", prettyURL, res.Status)
	}

	return resData, nil
}

// getUserProfileByID Get user profile by user ID
func getUserProfileByID(userID uuid.UUID) (*models.UserProfileModel, error) {
	profileURL := fmt.Sprintf("/profile/dto/id/%s", userID.String())
	foundProfileData, err := functionCall(http.MethodGet, []byte(""), profileURL, nil)
	if err != nil {
		if err == NotFoundHTTPStatusError {
			return nil, nil
		}
		log.Error("functionCall (%s) -  %s", profileURL, err.Error())
		return nil, fmt.Errorf("getUserProfileByID/functionCall")
	}
	var foundProfile models.UserProfileModel
	err = json.Unmarshal(foundProfileData, &foundProfile)
	if err != nil {
		log.Error("Unmarshal foundProfile -  %s", err.Error())
		return nil, fmt.Errorf("getUserProfileByID/unmarshal")
	}
	return &foundProfile, nil
}

// getCollectivesByID

func getCollectivesByID(collectiveID uuid.UUID) (*models.CollectivesModel, error) {
	collectiveURL := fmt.Sprintf("/collective/dto/id/%s", collectiveID.String())
	foundCollectivesData, err := functionCall(http.MethodGet, []byte(""), collectiveURL, nil)
	if err != nil {
		if err == NotFoundHTTPStatusError {
			return nil, nil
		}
		log.Error("functionCall (%s) -  %s", collectiveURL, err.Error())
		return nil, fmt.Errorf("getCollectivesByID/functionCall")
	}
	var foundCollectives models.CollectivesModel
	err = json.Unmarshal(foundCollectivesData, &foundCollectives)
	if err != nil {
		log.Error("Unmarshal foundCollectives -  %s", err.Error())
		return nil, fmt.Errorf("getCollectivesByID/unmarshal")
	}
	return &foundCollectives, nil
}

// readVotersAsync Read Voters async
func readVotersAsync(voterIds []string, infoReq *UserInfoInReq) <-chan ResultAsync {
	r := make(chan ResultAsync)
	go func() {
		defer close(r)
		votersURL := fmt.Sprintf("/users/?vin=%s", strings.Join(voterIds, ","))

		voters, err := functionCall(http.MethodGet, []byte(""), votersURL, getHeadersFromUserInfoReq(infoReq))
		if err != nil {
			r <- ResultAsync{Error: err}
			return
		}
		r <- ResultAsync{Result: voters}

	}()
	return r
}

func readVoters(c *fiber.Ctx, votes map[string]string) interface{} {
	userInfoReq := getUserInfoReq(c)

	voterIds := []string{}
	for voterId := range votes {
		voterIds = append(voterIds, voterId)
	}

	if len(voterIds) > 0 {
		readVotersChannel := readVotersAsync(voterIds, userInfoReq)

		votersResult := <-readVotersChannel
		if votersResult.Error != nil {
			messageError := fmt.Sprintf("Cannot get the voters! error: %s", votersResult.Error.Error())
			fmt.Println(messageError)
		} else {
			return votersResult.Result
		}
	}

	return votes
}
