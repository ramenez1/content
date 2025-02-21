package process

import (
	"Content/content_flow/internal/dao"
	"encoding/json"
	"fmt"
	flow "github.com/s8sg/goflow/flow/v1"
	goflow "github.com/s8sg/goflow/v1"
	"gorm.io/gorm"
	"log"
)

func ExecContentFlow(db *gorm.DB) {
	contentFlow := &contentFlow{
		contentDao: dao.NewContentDao(db),
	}
	fs := goflow.FlowService{
		Port:              7788,
		RedisURL:          "localhost:6379",
		WorkerConcurrency: 4,
	}
	_ = fs.Register("content-flow", contentFlow.flowHandle)
	err := fs.Start()
	if err != nil {
		fmt.Println("ExecContentFlow error : ", err)
		panic(err)
	}
}

func ExecContentWork(db *gorm.DB) {
	contentFlow := &contentFlow{
		contentDao: dao.NewContentDao(db),
	}
	fs := goflow.FlowService{
		Port:              7788,
		RedisURL:          "localhost:6379",
		WorkerConcurrency: 4,
	}
	_ = fs.Register("content-flow", contentFlow.flowHandle)
	err := fs.StartWorker()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func ExecServer(db *gorm.DB) {
	fs := goflow.FlowService{
		Port:              7788,
		RedisURL:          "localhost:6379",
		WorkerConcurrency: 5,
	}
	contentFlow := &contentFlow{
		contentDao: dao.NewContentDao(db),
	}
	_ = fs.Register("content-flow", contentFlow.flowHandle)
	err := fs.StartServer()
	if err != nil {
		panic(err)
	}
}

type contentFlow struct {
	contentDao *dao.ContentDao
}

func (c *contentFlow) flowHandle(workflow *flow.Workflow, context *flow.Context) error {
	dag := workflow.Dag()
	dag.Node("input", c.input)
	dag.Node("verify", c.verify)
	dag.Node("finish", c.finish)
	branches := dag.ConditionalBranch("branches",
		[]string{"category", "thumbnail", "pass", "format", "fail"},
		func(bytes []byte) []string {
			var data map[string]interface{}
			if err := json.Unmarshal(bytes, &data); err != nil {
				return nil
			}
			if data["approval_status"].(float64) == 2 {
				return []string{"category", "thumbnail", "pass", "format"}
			}
			return []string{"fail"}
		}, flow.Aggregator(func(m map[string][]byte) ([]byte, error) {
			fmt.Println(m)
			return []byte("ok"), nil
		}))
	branches["category"].Node("category", c.category)
	branches["thumbnail"].Node("thumbnail", c.thumbnail)
	branches["pass"].Node("category", c.pass)
	branches["format"].Node("format", c.format)
	branches["fail"].Node("fail", c.fail)

	dag.Edge("input", "verify")
	dag.Edge("verify", "branches")
	dag.Edge("branches", "finish")
	return nil
}

func (c *contentFlow) input(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec input")
	fmt.Println("input data : ", string(data))
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		log.Printf("Unmarshal error: %v", err)
		return nil, err
	}
	contentID := fmt.Sprintf("%v", input["content_id"])
	detail, err := c.contentDao.First(contentID)
	if err != nil {
		log.Printf("FirstFindByID err: %v", err)
		return nil, err
	}
	result, err := json.Marshal(map[string]interface{}{
		"title":      detail.Title,
		"video_url":  detail.VideoURL,
		"content_id": detail.ContentID,
	})
	if err != nil {
		log.Printf("MarshalResult err: %v", err)
		return nil, err
	}
	return result, nil
}

func (c *contentFlow) verify(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec verify")
	var detail map[string]interface{}
	if err := json.Unmarshal(data, &detail); err != nil {
		return nil, err
	}
	var (
		title     = detail["title"]
		videoURL  = detail["video_url"]
		contentID = detail["content_id"]
	)
	// 机审，人审
	detail["approval_status"] = 2
	fmt.Println(contentID, title, videoURL)
	return json.Marshal(detail)
}

func (c *contentFlow) category(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec category")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentID := input["content_id"].(string)
	err := c.contentDao.UpdateByID(contentID, "category", "category")
	if err != nil {
		return nil, err
	}
	return []byte("category"), nil
}

func (c *contentFlow) thumbnail(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec thumbnail")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentID := input["content_id"].(string)
	err := c.contentDao.UpdateByID(contentID, "thumbnail", "thumbnail")
	if err != nil {
		return nil, err
	}
	return []byte("thumbnail"), nil
}

func (c *contentFlow) format(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec format")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentID := input["content_id"].(string)
	err := c.contentDao.UpdateByID(contentID, "format", "format")
	if err != nil {
		return nil, err
	}
	return []byte("format"), nil
}

func (c *contentFlow) pass(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec pass")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentID := input["content_id"].(string)
	err := c.contentDao.UpdateByID(contentID, "approval_status", 2)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *contentFlow) fail(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec fail")
	var input map[string]interface{}
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	contentID := input["content_id"].(string)
	err := c.contentDao.UpdateByID(contentID, "approval_status", 3)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *contentFlow) finish(data []byte, option map[string][]string) ([]byte, error) {
	fmt.Println("exec finish")
	fmt.Println("result :", string(data))
	return data, nil
}
