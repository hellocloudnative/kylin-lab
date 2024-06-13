package vm

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	config2 "kylin-lab/tools/config"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type ServerWrapper struct {
	Servers []Server `json:"items"`
}

type Server struct {
	Status    string `json:"status"`
	Addresses struct {
		Vxlan []struct {
			OSEXTIPSMACMacAddr string `json:"OS-EXT-IPS-MAC:mac_addr"`
			Version            int    `json:"version"`
			Addr               string `json:"addr"`
			OSEXTIPSType       string `json:"OS-EXT-IPS:type"`
		} `json:"vxlan"`
	} `json:"addresses"`
	AvailabilityZone string `json:"availability_zone"`
	Image            struct {
		Id    string `json:"id"`
		Links []struct {
			Href string `json:"href"`
			Rel  string `json:"rel"`
		} `json:"links"`
	} `json:"image"`
	UpdatedAt  time.Time `json:"updated_at"`
	IsoMounted bool      `json:"iso_mounted"`
	Flavor     struct {
		Ephemeral    int    `json:"ephemeral"`
		Ram          int    `json:"ram"`
		OriginalName string `json:"original_name"`
		Vcpus        int    `json:"vcpus"`
		ExtraSpecs   struct {
		} `json:"extra_specs"`
		Swap int `json:"swap"`
		Disk int `json:"disk"`
	} `json:"flavor"`
	DeletedAt      interface{} `json:"deleted_at"`
	Id             string      `json:"id"`
	SecurityGroups []struct {
		Name string `json:"name"`
	} `json:"security_groups"`
	VolumesAttached []interface{} `json:"volumes_attached"`
	UserId          int           `json:"user_id"`
	BackendId       int           `json:"backend_id"`
	PowerState      int           `json:"power_state"`
	Metadata        struct {
	} `json:"metadata"`
	IsBaremetal bool        `json:"is_baremetal"`
	Description interface{} `json:"description"`
	Tags        string      `json:"tags"`
	Deleted     bool        `json:"deleted"`
	KeyName     interface{} `json:"key_name"`
	UserData    string      `json:"user_data"`
	Host        string      `json:"host"`
	RootBdm     struct {
		GuestFormat         interface{} `json:"guest_format"`
		BootIndex           int         `json:"boot_index"`
		AttachmentId        interface{} `json:"attachment_id"`
		DeleteOnTermination bool        `json:"delete_on_termination"`
		NoDevice            bool        `json:"no_device"`
		VolumeId            interface{} `json:"volume_id"`
		VolumeType          interface{} `json:"volume_type"`
		DeviceName          string      `json:"device_name"`
		DiskBus             interface{} `json:"disk_bus"`
		ImageId             string      `json:"image_id"`
		SourceType          string      `json:"source_type"`
		Tag                 interface{} `json:"tag"`
		DeviceType          string      `json:"device_type"`
		SnapshotId          interface{} `json:"snapshot_id"`
		DestinationType     string      `json:"destination_type"`
		VolumeSize          interface{} `json:"volume_size"`
	} `json:"root_bdm"`
	TaskState interface{} `json:"task_state"`
	Locked    bool        `json:"locked"`
	Name      string      `json:"name"`
	TenantId  string      `json:"tenant_id"`
	CreatedAt time.Time   `json:"created_at"`
}

type ServerInfoResponse struct {
	Servers []struct {
		Status    string `json:"status"`
		Addresses struct {
		} `json:"addresses"`
		AvailabilityZone string `json:"availability_zone"`
		Image            struct {
			Id    string `json:"id"`
			Links []struct {
				Href string `json:"href"`
				Rel  string `json:"rel"`
			} `json:"links"`
		} `json:"image"`
		UpdatedAt  time.Time `json:"updated_at"`
		IsoMounted bool      `json:"iso_mounted"`
		Flavor     struct {
			Ephemeral    int    `json:"ephemeral"`
			Ram          int    `json:"ram"`
			OriginalName string `json:"original_name"`
			Vcpus        int    `json:"vcpus"`
			ExtraSpecs   struct {
			} `json:"extra_specs"`
			Swap int `json:"swap"`
			Disk int `json:"disk"`
		} `json:"flavor"`
		DeletedAt      interface{} `json:"deleted_at"`
		Id             string      `json:"id"`
		SecurityGroups []struct {
			Name string `json:"name"`
		} `json:"security_groups"`
		VolumesAttached []interface{} `json:"volumes_attached"`
		UserId          int           `json:"user_id"`
		BackendId       int           `json:"backend_id"`
		PowerState      int           `json:"power_state"`
		Metadata        struct {
		} `json:"metadata"`
		IsBaremetal bool          `json:"is_baremetal"`
		Description interface{}   `json:"description"`
		Tags        []interface{} `json:"tags"`
		Deleted     bool          `json:"deleted"`
		KeyName     interface{}   `json:"key_name"`
		UserData    string        `json:"user_data"`
		Host        interface{}   `json:"host"`
		RootBdm     struct {
		} `json:"root_bdm"`
		TaskState string    `json:"task_state"`
		Locked    bool      `json:"locked"`
		Name      string    `json:"name"`
		TenantId  string    `json:"tenant_id"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"servers"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type ImageWrapper struct {
	Images []Image `json:"images"`
}

type FlavorsWrapper struct {
	Flavors []Flavor `json:"flavors"`
}

type VNCWrapper struct {
	Passwd   interface{} `json:"passwd"`
	Host     string      `json:"host"`
	Protocol string      `json:"protocol"`
	Url      string      `json:"url"`
	Tlsport  interface{} `json:"tlsport"`
	Type     string      `json:"type"`
	Port     string      `json:"port"`
}

type NetworkWrapper struct {
	Networks []Network `json:"networks"`
}

type Image struct {
	Status           string      `json:"status"`
	OsDistro         string      `json:"os_distro"`
	UserId           interface{} `json:"user_id"`
	Name             string      `json:"name"`
	Tags             string      `json:"tags"`
	ImageType        string      `json:"image_type"`
	Checksum         string      `json:"checksum"`
	MinRam           int         `json:"min_ram"`
	Description      interface{} `json:"description"`
	DiskFormat       string      `json:"disk_format"`
	Visibility       string      `json:"visibility"`
	Baremetal        bool        `json:"baremetal"`
	Owner            string      `json:"owner"`
	HwQemuGuestAgent string      `json:"hw_qemu_guest_agent"`
	InstanceUuid     interface{} `json:"instance_uuid"`
	MinDisk          int         `json:"min_disk"`
	CreatedAt        time.Time   `json:"created_at"`
	VirtualSize      int64       `json:"virtual_size"`
	HwArchitecture   string      `json:"hw_architecture"`
	Id               string      `json:"id"`
	Size             int64       `json:"size"`
}

type Flavor struct {
	UserId      interface{} `json:"user_id"`
	Name        string      `json:"name"`
	Ram         int         `json:"ram"`
	Ephemeral   int         `json:"ephemeral"`
	Vcpus       int         `json:"vcpus"`
	Swap        int         `json:"swap"`
	RxtxFactor  float64     `json:"rxtx_factor"`
	IsPublic    bool        `json:"is_public"`
	Disk        int         `json:"disk"`
	Id          string      `json:"id"`
	Description interface{} `json:"description"`
}

type Network struct {
	Status            string      `json:"status"`
	Subnets           []string    `json:"subnets"`
	AvailabilityZones []string    `json:"availability_zones"`
	UserId            int         `json:"user_id"`
	Name              string      `json:"name"`
	AdminStateUp      bool        `json:"admin_state_up"`
	SegmentationId    int         `json:"segmentation_id"`
	CreatedAt         time.Time   `json:"created_at"`
	Tags              string      `json:"tags"`
	Id                string      `json:"id"`
	Description       string      `json:"description"`
	PhysicalNetwork   interface{} `json:"physical_network"`
	External          bool        `json:"external"`
	TenantId          string      `json:"tenant_id"`
	Shared            bool        `json:"shared"`
	Mtu               int         `json:"mtu"`
	NetworkType       string      `json:"network_type"`
	QosPolicyId       interface{} `json:"qos_policy_id"`
}

type KylinCloudInstance struct {
	Name       string `json:"name"`
	SourceId   string `json:"source_id"`
	FlavorId   string `json:"flavor_id"`
	SourceType string `json:"source_type"`
	Nics       []struct {
		NetId string `json:"net-id"`
	} `json:"nics"`
	UserData       string   `json:"user_data"`
	SecurityGroups []string `json:"security_groups"`
}

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	for i, v := range b {
		// 将随机数转换为字符集中的字符
		b[i] = charset[v%byte(len(charset))]
	}
	return string(b), nil
}

func PostAuthRequestToken(urlsting string) (string, error) {
	var authRequest = AuthRequest{
		Username: "system",
		Password: "sys@cloud_",
	}

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// 将AuthRequest转换为form-data
	formData := url.Values{}
	formData.Set("username", authRequest.Username)
	formData.Set("password", authRequest.Password)

	// 创建HTTP POST请求
	req, err := http.NewRequest(http.MethodPost, urlsting, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return "", err
	}

	// 设置请求头部，指定发送的数据是form-data格式
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to authenticate: %v", resp.StatusCode)
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return "", err
	}

	// 解析响应体，提取token
	var authResponse AuthResponse
	if err := json.Unmarshal(body, &authResponse); err != nil {
		log.Println("Error unmarshalling auth response:", err)
		return "", err
	}

	// 返回获取到的token
	return authResponse.Data.Token, nil
}

func GetImagesRequest() (ImageWrapper, error) {
	// 调用 PostAuthRequestToken 函数获取 token
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		return ImageWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 创建HTTP GET请求
	req, err := http.NewRequest(http.MethodGet, config2.KylinCloudConfig.ApiUrl+"/api/images", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return ImageWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 设置请求头部
	req.Header.Set("x-auth-token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return ImageWrapper{}, err // 返回空的 ImageWrapper 和错误
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get images: %v\n", resp.StatusCode)
		return ImageWrapper{}, fmt.Errorf("failed to get images: %v", resp.StatusCode) // 返回空的 ImageWrapper 和错误
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ImageWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 解析响应体为 ImageWrapper
	var imagesResponse ImageWrapper
	if err := json.Unmarshal(body, &imagesResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		return ImageWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	//log.Println("Parsed images:", imagesResponse.Images)
	return imagesResponse, nil // 成功返回解析后的 ImageWrapper 和 nil
}

func GetServersRequest(id string) (ServerWrapper, error) {
	// 调用 PostAuthRequestToken 函数获取 token
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		return ServerWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	baseURL := config2.KylinCloudConfig.ApiUrl + "/api/instances"
	queryParams := url.Values{}
	queryParams.Add("id", id) // 添加查询参数 id=123
	urlWithQuery := baseURL + "?" + queryParams.Encode()
	// 创建HTTP GET请求
	req, err := http.NewRequest(http.MethodGet, urlWithQuery, nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return ServerWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 设置请求头部
	req.Header.Set("x-auth-token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return ServerWrapper{}, err // 返回空的 ImageWrapper 和错误
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get images: %v\n", resp.StatusCode)
		return ServerWrapper{}, fmt.Errorf("failed to get server: %v", resp.StatusCode) // 返回空的 ImageWrapper 和错误
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return ServerWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 解析响应体为 ImageWrapper
	var serverssResponse ServerWrapper
	if err := json.Unmarshal(body, &serverssResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		return ServerWrapper{}, err
	}

	return serverssResponse, nil
}

func GetVNCRequest(id string) (VNCWrapper, error) {
	// 调用 PostAuthRequestToken 函数获取 token
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		return VNCWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 创建HTTP GET请求
	req, err := http.NewRequest(http.MethodGet, config2.KylinCloudConfig.ApiUrl+"/api/instances/"+id+"/console", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return VNCWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 设置请求头部
	req.Header.Set("x-auth-token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return VNCWrapper{}, err // 返回空的 ImageWrapper 和错误
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get images: %v\n", resp.StatusCode)
		return VNCWrapper{}, fmt.Errorf("failed to get images: %v", resp.StatusCode) // 返回空的 ImageWrapper 和错误
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return VNCWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 解析响应体为 ImageWrapper
	var vncResponse VNCWrapper
	if err := json.Unmarshal(body, &vncResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		return VNCWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	//log.Println("Parsed images:", imagesResponse.Images)
	return vncResponse, nil // 成功返回解析后的 ImageWrapper 和 nil
}

func GetFlavorsRequest() (FlavorsWrapper, error) {
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		return FlavorsWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 创建HTTP GET请求
	req, err := http.NewRequest(http.MethodGet, config2.KylinCloudConfig.ApiUrl+"/api/flavors", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return FlavorsWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 设置请求头部
	req.Header.Set("x-auth-token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return FlavorsWrapper{}, err // 返回空的 ImageWrapper 和错误
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get images: %v\n", resp.StatusCode)
		return FlavorsWrapper{}, fmt.Errorf("failed to get flavors: %v", resp.StatusCode) // 返回空的 ImageWrapper 和错误
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return FlavorsWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 解析响应体为 ImageWrapper
	var flavorsrResponse FlavorsWrapper
	if err := json.Unmarshal(body, &flavorsrResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		return FlavorsWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	//log.Println("Parsed images:", flavorsrResponse.Flavors)
	return flavorsrResponse, nil //
}

func GetNetworksRequest() (NetworkWrapper, error) {
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		return NetworkWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 创建HTTP GET请求
	req, err := http.NewRequest(http.MethodGet, config2.KylinCloudConfig.ApiUrl+"/api/networks", nil)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return NetworkWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 设置请求头部
	req.Header.Set("x-auth-token", token)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Second * 10, // 设置超时时间
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending GET request:", err)
		return NetworkWrapper{}, err // 返回空的 ImageWrapper 和错误
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get images: %v\n", resp.StatusCode)
		return NetworkWrapper{}, fmt.Errorf("failed to get images: %v", resp.StatusCode) // 返回空的 ImageWrapper 和错误
	}

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return NetworkWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	// 解析响应体为 ImageWrapper
	var networkResponse NetworkWrapper
	if err := json.Unmarshal(body, &networkResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		return NetworkWrapper{}, err // 返回空的 ImageWrapper 和错误
	}

	//log.Println("Parsed images:", networkResponse.Networks)
	return networkResponse, nil //
}

func PostApplyInstances(hwArchitecture, imageName, Flavors, network string) (ServerInfoResponse, error) {
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		return ServerInfoResponse{}, err
	}

	randomStringLength := 5
	randomString, err := generateRandomString(randomStringLength)
	if err != nil {
		log.Error("Error generating random string:", err)
		return ServerInfoResponse{}, err
	}
	if hwArchitecture != GetHwArchitecture(imageName) {
		log.Error("Hardware architecture does not match the image.")
		return ServerInfoResponse{}, err
	}
	var kylinCloudInstance KylinCloudInstance
	kylinCloudInstance.Name = "instance-" + randomString
	kylinCloudInstance.SourceId = GetSourceId(imageName) // 假设这是从某个函数获取的
	kylinCloudInstance.SourceType = "image"
	kylinCloudInstance.FlavorId = GetFlavorId(Flavors) // 同上
	kylinCloudInstance.Nics = append(kylinCloudInstance.Nics, struct {
		NetId string `json:"net-id"`
	}{NetId: GetNetworkId(network)})

	kylinCloudInstance.UserData = "ga"
	kylinCloudInstance.SecurityGroups = []string{}

	log.Info("KylinCloudInstance:", kylinCloudInstance)

	jsonBody, err := json.Marshal(kylinCloudInstance)
	if err != nil {
		log.Error("Error marshalling KylinCloudInstance to JSON:", err)
		return ServerInfoResponse{}, err
	}

	req, err := http.NewRequest(http.MethodPost, config2.KylinCloudConfig.ApiUrl+"/api/instances/", bytes.NewReader(jsonBody))
	if err != nil {
		log.Error("Error creating HTTP request:", err)
		return ServerInfoResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("x-auth-token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending POST request:", err)
		return ServerInfoResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("Failed to create instance: %v, Response: %s\n", resp.StatusCode, body)
		return ServerInfoResponse{}, fmt.Errorf("failed to create instance, status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return ServerInfoResponse{}, err
	}

	// 解析响应体为 ImageWrapper
	var serverInfoResponse ServerInfoResponse
	if err := json.Unmarshal(body, &serverInfoResponse); err != nil {
		log.Println("Error unmarshalling response:", err)
		return ServerInfoResponse{}, err // 返回空的 ImageWrapper 和错误
	}

	//log.Println("Parsed images:", imagesResponse.Images)
	return serverInfoResponse, nil // 成功返回解析后的 ImageWrapper 和 nil

}

func GetSourceId(imageName string) string {
	images, err := GetImagesRequest()
	if err != nil {
		log.Error(err)
	}
	for image := range images.Images {
		if images.Images[image].Name == imageName {
			return images.Images[image].Id
		}
	}

	log.Error("Image not found")
	return ""
}

func GetHwArchitecture(imageName string) string {
	images, err := GetImagesRequest()
	if err != nil {
		log.Error(err)
		return ""
	}

	// 检查 images.Images 是否为空
	if len(images.Images) > 0 {
		// 使用 for-range 循环迭代 images.Images 切片
		for _, image := range images.Images {
			if image.Name == imageName {
				return image.HwArchitecture
			}
		}
	}

	// 如果没有找到匹配的 Image，记录错误
	log.Error("Image not found")
	return ""
}

func GetFlavorId(flavorsName string) string {
	flavorsResponse, err := GetFlavorsRequest()
	if err != nil {
		log.Error(err)
	}
	// 假设 flavorsResponse.Flavors 是一个 map 类型的字段
	for flavor := range flavorsResponse.Flavors {
		if flavorsResponse.Flavors[flavor].Name == convertToKylinCloudFormat(flavorsName) {
			return flavorsResponse.Flavors[flavor].Id
		}
	}

	// 如果没有找到匹配的 Flavor，记录错误
	log.Error("Flavor not found")
	log.Info("ASAAAAAAA", convertToKylinCloudFormat(flavorsName))
	return ""
}

func GetNetworkId(networkName string) string {
	networks, err := GetNetworksRequest()
	if err != nil {
		log.Error(err)
	}
	for network := range networks.Networks {
		if networks.Networks[network].Name == networkName {
			return networks.Networks[network].Id
		}
	}
	log.Error("Network not found")
	return ""
}

func convertToKylinCloudFormat(input string) string {
	replacement := "c."
	return "kylincloud." + strings.Replace(input, "C-", replacement, 1)
}

func GetKylinCloudToken(c *gin.Context) {
	// 调用函数从指定的HTTPS URL获取token
	token, err := PostAuthRequestToken(config2.KylinCloudConfig.AuthUrl + "/api/auth")
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token"})
		return
	}

	// 将获取到的token返回给客户端
	c.JSON(http.StatusOK, gin.H{"data": gin.H{"token": token}})
}

func GetKylinCloudImages(c *gin.Context) {
	// 调用函数从指定的HTTPS URL获取token
	images, err := GetImagesRequest()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve token"})
		return
	}

	c.JSON(http.StatusOK, images)
}
