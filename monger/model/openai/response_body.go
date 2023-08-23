package model

/*
	{
	  "id": "chatcmpl-7iyWV166Cahu7A08Iw8FrG3M4LXnd",
	  "object": "chat.completion",
	  "created": 1690953439,
	  "model": "gpt-3.5-turbo-0613",
	  "choices": [
	    {
	      "index": 0,
	      "message": {
	        "role": "assistant",
	        "content": "I'm sorry, but I am an AI language model and I am unable to physically feed you anything. However, I can provide you with a recipe for chocolate chip cookies if you'd like!"
	      },
	      "finish_reason": "stop"
	    }
	  ],
	  "usage": {
	    "prompt_tokens": 12,
	    "completion_tokens": 39,
	    "total_tokens": 51
	  }
	}
*/

type ResponseBody struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
	Usage   Usage     `json:"usage,omitempty"`
}
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type Choices struct {
	Index        int     `json:"index"`
	Message      Message `json:"message,omitempty"`
	Delta        Delta   `json:"delta,omitempty"`
	FinishReason string  `json:"finish_reason"`
}
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}
type Delta struct {
	Content string `json:"content"`
}

type Content struct {
	Content string `json:"content"`
}
