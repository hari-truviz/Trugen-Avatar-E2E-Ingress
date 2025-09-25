package CustomTypes

type ConversationRequest struct {
	Input Input `json:"input"`
}

type Input struct {
	Avatar Avatar `json:"avatar"`
}

type Avatar struct {
	Timeout                int                 `json:"timeout"`
	UserID                 string              `json:"user_id"`
	AgentID                string              `json:"agent_id"`
	ConversationID         string              `json:"conversation_id"`
	RoomID                 string              `json:"room_id"`
	AvatarID               string              `json:"avatar_id"`
	AvatarDataSource       string              `json:"avatar_data_source"`
	FrameRate              int                 `json:"frame_rate"`
	MaxMessageHistory      int                 `json:"max_message_history"`
	SilencePadding         int                 `json:"silence_padding"`
	IsFaceEnhancerEnabled  bool                `json:"is_face_enhancer_enabled"`
	FaceEnhancerModelName  string              `json:"face_enhancer_model_name"`
	FaceEnhancerDenoise    float64             `json:"face_enhancer_denoise_strength"`
	PersonaName            string              `json:"persona_name"`
	PersonaPrompt          string              `json:"persona_prompt"`
	ConversationalContext  string              `json:"conversational_context"`
	Memory                 Memory              `json:"memory"`
	InterpolationConfig    InterpolationConfig `json:"interpolation_config"`
	SceneContextEngine     SceneContextEngine  `json:"scene_context_engine"`
	IdleTimeout            IdleTimeout         `json:"idle_timeout"`
	SceneAnalyzerPrompt    SceneAnalyzerPrompt `json:"scene_analyzer_prompt"`
	WarningExitMessage     ExitMessage         `json:"warning_exit_message"`
	ExitMessage            ExitMessage         `json:"exit_message"`
	WelcomeMessage         WelcomeMessage      `json:"welcome_message"`
	EyeMaskReplacement     bool                `json:"eye_mask_replacement"`
	AudioFeaturesType      string              `json:"audio_features_type"`
	AudioFeaturesWindowLen int                 `json:"audio_features_window_length"`
	PrevAudioDuration      int                 `json:"prev_audio_duration"`
	Config                 Config              `json:"config"`
	KnowledgeBase          interface{}         `json:"knowledge_base"`
	McpServers             []McpServer         `json:"mcp_servers"`
	Tools                  []Tool              `json:"tools"`
}

type Memory struct {
	Enabled  bool   `json:"enabled"`
	Excludes string `json:"excludes"`
	Includes string `json:"includes"`
}

type InterpolationConfig struct {
	Enabled bool `json:"enabled"`
	Exp     int  `json:"exp"`
}

type SceneContextEngine struct {
	OnSnapshotTimeout int        `json:"on_snapshot_timeout"`
	SnapshotScale     float64    `json:"snapshot_scale"`
	VisionLLM         string     `json:"vision_llm"`
	LLMPrompts        LLMPrompts `json:"llm_prompts"`
}

type LLMPrompts struct {
	GetUserAppearance                    string       `json:"get_user_appearance"`
	FirstQuery                           string       `json:"first_query"`
	AnalyzeSceneCtxResponse              string       `json:"analyze_scene_ctx_response"`
	AnalyzeActionsSystemPrompt           string       `json:"analyze_actions_system_prompt"`
	AnalyzeAction                        string       `json:"analyze_action"`
	ActionsList                          []ActionItem `json:"actions_list"`
	AddActionRecognitionSyntheticUserQry bool         `json:"add_action_recognition_synthetic_user_query"`
	SyntheticUserQuery                   string       `json:"synthetic_user_query"`
	UserQueryAnalysisSystemPrompt        string       `json:"user_query_analysis_system_prompt"`
}

type ActionItem struct {
	ActionName              string `json:"Action_Name"`
	Type                    string `json:"Type"`
	AnalysisInstruction     string `json:"Analysis_Instruction"`
	ActionNeedsToBeObserved string `json:"Action_Needs_To_Be_Observed"`
}

type IdleTimeout struct {
	Timeout       int      `json:"timeout"`
	FillerPhrases []string `json:"filler_phrases"`
}

type SceneAnalyzerPrompt struct {
	SystemPrompt string `json:"system_prompt"`
	TaskPrompt   string `json:"task_prompt"`
}

type ExitMessage struct {
	CalloutBefore int      `json:"callout_before,omitempty"`
	MaxCallDur    int      `json:"max_call_duration,omitempty"`
	Messages      []string `json:"messages"`
}

type WelcomeMessage struct {
	WaitTime int      `json:"wait_time"`
	Messages []string `json:"messages"`
}

type Config struct {
	PreemptiveSynthesis bool            `json:"preemptive_synthesis"`
	Protocol            Protocol        `json:"protocol"`
	Snapshot            Snapshot        `json:"snapshot"`
	NoiseCancellation   NoiseCancel     `json:"noise_cancellation"`
	SuperResolution     SuperResolution `json:"super_resolution"`
	VAD                 VAD             `json:"vad"`
	STT                 STT             `json:"stt"`
	TurnDetector        bool            `json:"turn_detector"`
	LLM                 LLM             `json:"llm"`
	TTS                 TTS             `json:"tts"`
}

type Protocol struct {
	VideoCodec   string `json:"video_codec"`
	VideoBitrate int    `json:"video_bitrate"`
	Simulcast    bool   `json:"simulcast"`
}

type Snapshot struct {
	Enabled bool    `json:"enabled"`
	Scale   float64 `json:"scale"`
}

type NoiseCancel struct {
	Provider string `json:"provider"`
}

type SuperResolution struct {
	Enabled bool `json:"enabled"`
	Scale   int  `json:"scale"`
}

type VAD struct {
	Enabled             bool    `json:"enabled"`
	Provider            string  `json:"provider"`
	MinSilenceDuration  float64 `json:"min_silence_duration"`
	ActivationThreshold float64 `json:"activation_threshold"`
}

type STT struct {
	Provider                       string  `json:"provider"`
	Model                          string  `json:"model"`
	Language                       string  `json:"language"`
	FallbackModel                  string  `json:"fallback_model"`
	InterruptSpeechDuration        float64 `json:"interrupt_speech_duration"`
	AllowIntermResultsInterruption bool    `json:"allow_interm_results_interruption"`
	MinEndpointingDelay            float64 `json:"min_endpointing_delay"`
	MaxEndpointingDelay            float64 `json:"max_endpointing_delay"`
	FinalTranscriptTimeout         float64 `json:"final_transcript_timeout"`
	PunctuationReduceFactor        float64 `json:"punctuation_reduce_factor"`
}

type LLM struct {
	Provider        string      `json:"provider"`
	Model           string      `json:"model"`
	FallbackModel   string      `json:"fallback_model"`
	UseNLTK         bool        `json:"use_nltk"`
	ReasoningEffort interface{} `json:"reasoning_effort"`
}

type TTS struct {
	Provider         string  `json:"provider"`
	ModelID          string  `json:"model_id"`
	Language         string  `json:"language"`
	VoiceID          string  `json:"voice_id"`
	Pitch            int     `json:"pitch"`
	EffectsProfileID string  `json:"effects_profile_id"`
	SpeakingRate     int     `json:"speaking_rate"`
	Stability        float32 `json:"stability"`
	SimilarityBoost  float32 `json:"similarity_boost"`
	SampleRate       int     `json:"sample_rate"`
	Encoding         string  `json:"encoding"`
	Gender           string  `json:"gender"`
	FallbackVoiceID  string  `json:"fallback_voice_id"`
}

//
// MCP servers and tools (added concrete structs to match JSON shape)
//

type McpServer struct {
	Type           string                 `json:"type"`
	Name           string                 `json:"name"`
	Params         map[string]interface{} `json:"params"`
	CacheToolsList bool                   `json:"cache_tools_list"`
	EventMessages  EventMessages          `json:"event_messages"`
}

type EventMessages struct {
	OnStart   *SimpleMessage `json:"on_start,omitempty"`
	OnSuccess *SimpleMessage `json:"on_success,omitempty"`
	OnDelay   *DelayMessage  `json:"on_delay,omitempty"`
	OnError   *SimpleMessage `json:"on_error,omitempty"`
}

type SimpleMessage struct {
	Message string `json:"message"`
}

type DelayMessage struct {
	Delay   int    `json:"delay"`
	Message string `json:"message"`
}

//
// Tools
//

type Tool struct {
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Arguments     []ToolArgument    `json:"arguments"`
	RequestConfig RequestConfig     `json:"request_config"`
	EventMessages ToolEventMessages `json:"event_messages"`
}

type ToolArgument struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type RequestConfig struct {
	Method  string            `json:"method"`
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers"`
}

type ToolEventMessages struct {
	OnStart   *SimpleMessage `json:"on_start,omitempty"`
	OnDelay   *DelayMessage  `json:"on_delay,omitempty"`
	OnSuccess *SimpleMessage `json:"on_success,omitempty"`
	OnError   *SimpleMessage `json:"on_error,omitempty"`
}

type Acknowledgment struct {
	Message         string `json:"message"`
	CurrentPosition int    `json:"position"`
	StatusCode      int    `json:"status"`
}

type AllAvatars struct {
	Avatars []Avatar `json:"avatars"`
}
