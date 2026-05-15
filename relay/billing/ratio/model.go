package ratio

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/songquanpeng/one-api/common/logger"
)

const (
	USD2RMB = 7
	USD     = 500 // $0.002 = 1 -> $1 = 500
	RMB     = USD / USD2RMB
)

// ModelRatio
// https://platform.openai.com/docs/models/model-endpoint-compatibility
// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Blfmc9dlf
// https://openai.com/pricing
// 1 === $0.002 / 1K tokens
// 1 === ￥0.014 / 1k tokens
var ModelRatio = map[string]ModelRatioDetail{
	// https://openai.com/pricing
	"gpt-4":                   {NormalRatio: 15, CachedRatio: 15},
	"gpt-4-0314":              {NormalRatio: 15, CachedRatio: 15},
	"gpt-4-0613":              {NormalRatio: 15, CachedRatio: 15},
	"gpt-4-32k":               {NormalRatio: 30, CachedRatio: 30},
	"gpt-4-32k-0314":          {NormalRatio: 30, CachedRatio: 30},
	"gpt-4-32k-0613":          {NormalRatio: 30, CachedRatio: 30},
	"gpt-4-1106-preview":      {NormalRatio: 5, CachedRatio: 5},
	"gpt-4-0125-preview":      {NormalRatio: 5, CachedRatio: 5},
	"gpt-4-turbo-preview":     {NormalRatio: 5, CachedRatio: 5},
	"gpt-4-turbo":             {NormalRatio: 5, CachedRatio: 5},
	"gpt-4-turbo-2024-04-09":  {NormalRatio: 5, CachedRatio: 5},
	"gpt-4o":                  {NormalRatio: 2.5, CachedRatio: 2.5},     // $0.005 / 1K tokens
	"chatgpt-4o-latest":       {NormalRatio: 2.5, CachedRatio: 2.5},     // $0.005 / 1K tokens
	"gpt-4o-2024-05-13":       {NormalRatio: 2.5, CachedRatio: 2.5},     // $0.005 / 1K tokens
	"gpt-4o-2024-08-06":       {NormalRatio: 1.25, CachedRatio: 1.25},   // $0.0025 / 1K tokens
	"gpt-4o-mini":             {NormalRatio: 0.075, CachedRatio: 0.075}, // $0.00015 / 1K tokens
	"gpt-4o-mini-2024-07-18":  {NormalRatio: 0.075, CachedRatio: 0.075}, // $0.00015 / 1K tokens
	"gpt-5-chat":              {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5":                   {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5-2025-08-07":        {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5-chat-latest":       {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5-mini":              {NormalRatio: 0.125, CachedRatio: 0.125},
	"gpt-5-mini-2025-08-07":   {NormalRatio: 0.125, CachedRatio: 0.125},
	"gpt-5.4":                 {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.4-mini":            {NormalRatio: 0.125, CachedRatio: 0.125},
	"gpt-5.4-pro":             {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.4-nano":            {NormalRatio: 0.025, CachedRatio: 0.025},
	"gpt-5.5":                 {NormalRatio: 0.88, CachedRatio: 0.09},
	"gpt-5-nano":              {NormalRatio: 0.025, CachedRatio: 0.025},
	"gpt-5.1":                 {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.1-chat":            {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.1-codex":           {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.1-codex-mini":      {NormalRatio: 0.125, CachedRatio: 0.125},
	"gpt-5-codex":             {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.2":                 {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-5.2-codex":           {NormalRatio: 0.625, CachedRatio: 0.625},
	"gpt-4-vision-preview":    {NormalRatio: 5, CachedRatio: 5},       // $0.01 / 1K tokens
	"gpt-3.5-turbo":           {NormalRatio: 0.25, CachedRatio: 0.25}, // $0.0005 / 1K tokens
	"gpt-3.5-turbo-0301":      {NormalRatio: 0.75, CachedRatio: 0.75},
	"gpt-3.5-turbo-0613":      {NormalRatio: 0.75, CachedRatio: 0.75},
	"gpt-3.5-turbo-16k":       {NormalRatio: 1.5, CachedRatio: 1.5}, // $0.003 / 1K tokens
	"gpt-3.5-turbo-16k-0613":  {NormalRatio: 1.5, CachedRatio: 1.5},
	"gpt-3.5-turbo-instruct":  {NormalRatio: 0.75, CachedRatio: 0.75}, // $0.0015 / 1K tokens
	"gpt-3.5-turbo-1106":      {NormalRatio: 0.5, CachedRatio: 0.5},   // $0.001 / 1K tokens
	"gpt-3.5-turbo-0125":      {NormalRatio: 0.25, CachedRatio: 0.25}, // $0.0005 / 1K tokens
	"davinci-002":             {NormalRatio: 1, CachedRatio: 1},       // $0.002 / 1K tokens
	"babbage-002":             {NormalRatio: 0.2, CachedRatio: 0.2},   // $0.0004 / 1K tokens
	"text-ada-001":            {NormalRatio: 0.2, CachedRatio: 0.2},
	"text-babbage-001":        {NormalRatio: 0.25, CachedRatio: 0.25},
	"text-curie-001":          {NormalRatio: 1, CachedRatio: 1},
	"text-davinci-002":        {NormalRatio: 10, CachedRatio: 10},
	"text-davinci-003":        {NormalRatio: 10, CachedRatio: 10},
	"text-davinci-edit-001":   {NormalRatio: 10, CachedRatio: 10},
	"code-davinci-edit-001":   {NormalRatio: 10, CachedRatio: 10},
	"whisper-1":               {NormalRatio: 15, CachedRatio: 15},   // $0.006 / minute -> $0.006 / 150 words -> $0.006 / 200 tokens -> $0.03 / 1k tokens
	"tts-1":                   {NormalRatio: 7.5, CachedRatio: 7.5}, // $0.015 / 1K characters
	"tts-1-1106":              {NormalRatio: 7.5, CachedRatio: 7.5},
	"tts-1-hd":                {NormalRatio: 15, CachedRatio: 15}, // $0.030 / 1K characters
	"tts-1-hd-1106":           {NormalRatio: 15, CachedRatio: 15},
	"davinci":                 {NormalRatio: 10, CachedRatio: 10},
	"curie":                   {NormalRatio: 10, CachedRatio: 10},
	"babbage":                 {NormalRatio: 10, CachedRatio: 10},
	"ada":                     {NormalRatio: 10, CachedRatio: 10},
	"text-embedding-ada-002":  {NormalRatio: 0.05, CachedRatio: 0.05},
	"text-embedding-3-small":  {NormalRatio: 0.01, CachedRatio: 0.01},
	"text-embedding-3-large":  {NormalRatio: 0.065, CachedRatio: 0.065},
	"text-search-ada-doc-001": {NormalRatio: 10, CachedRatio: 10},
	"text-moderation-stable":  {NormalRatio: 0.1, CachedRatio: 0.1},
	"text-moderation-latest":  {NormalRatio: 0.1, CachedRatio: 0.1},
	"dall-e-2":                {NormalRatio: 0.02 * USD, CachedRatio: 0.02 * USD}, // $0.016 - $0.020 / image
	"dall-e-3":                {NormalRatio: 0.04 * USD, CachedRatio: 0.04 * USD}, // $0.040 - $0.120 / image
	// https://www.anthropic.com/api#pricing
	"claude-instant-1.2":         {NormalRatio: 0.8 / 1000 * USD, CachedRatio: 0.8 / 1000 * USD},
	"claude-2.0":                 {NormalRatio: 8.0 / 1000 * USD, CachedRatio: 8.0 / 1000 * USD},
	"claude-2.1":                 {NormalRatio: 8.0 / 1000 * USD, CachedRatio: 8.0 / 1000 * USD},
	"claude-3-haiku-20240307":    {NormalRatio: 0.25 / 1000 * USD, CachedRatio: 0.25 / 1000 * USD},
	"claude-3-5-haiku-20241022":  {NormalRatio: 1.0 / 1000 * USD, CachedRatio: 1.0 / 1000 * USD},
	"claude-3-sonnet-20240229":   {NormalRatio: 3.0 / 1000 * USD, CachedRatio: 3.0 / 1000 * USD},
	"claude-3-5-sonnet-20240620": {NormalRatio: 3.0 / 1000 * USD, CachedRatio: 3.0 / 1000 * USD},
	"claude-3-5-sonnet-20241022": {NormalRatio: 3.0 / 1000 * USD, CachedRatio: 3.0 / 1000 * USD},
	"claude-3-opus-20240229":     {NormalRatio: 15.0 / 1000 * USD, CachedRatio: 15.0 / 1000 * USD},
	// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/hlrk4akp7
	"ERNIE-4.0-8K":       {NormalRatio: 0.120 * RMB, CachedRatio: 0.120 * RMB},
	"ERNIE-3.5-8K":       {NormalRatio: 0.012 * RMB, CachedRatio: 0.012 * RMB},
	"ERNIE-3.5-8K-0205":  {NormalRatio: 0.024 * RMB, CachedRatio: 0.024 * RMB},
	"ERNIE-3.5-8K-1222":  {NormalRatio: 0.012 * RMB, CachedRatio: 0.012 * RMB},
	"ERNIE-Bot-8K":       {NormalRatio: 0.024 * RMB, CachedRatio: 0.024 * RMB},
	"ERNIE-3.5-4K-0205":  {NormalRatio: 0.012 * RMB, CachedRatio: 0.012 * RMB},
	"ERNIE-Speed-8K":     {NormalRatio: 0.004 * RMB, CachedRatio: 0.004 * RMB},
	"ERNIE-Speed-128K":   {NormalRatio: 0.004 * RMB, CachedRatio: 0.004 * RMB},
	"ERNIE-Lite-8K-0922": {NormalRatio: 0.008 * RMB, CachedRatio: 0.008 * RMB},
	"ERNIE-Lite-8K-0308": {NormalRatio: 0.003 * RMB, CachedRatio: 0.003 * RMB},
	"ERNIE-Tiny-8K":      {NormalRatio: 0.001 * RMB, CachedRatio: 0.001 * RMB},
	"BLOOMZ-7B":          {NormalRatio: 0.004 * RMB, CachedRatio: 0.004 * RMB},
	"Embedding-V1":       {NormalRatio: 0.002 * RMB, CachedRatio: 0.002 * RMB},
	"bge-large-zh":       {NormalRatio: 0.002 * RMB, CachedRatio: 0.002 * RMB},
	"bge-large-en":       {NormalRatio: 0.002 * RMB, CachedRatio: 0.002 * RMB},
	"tao-8k":             {NormalRatio: 0.002 * RMB, CachedRatio: 0.002 * RMB},
	// https://ai.google.dev/pricing
	"gemini-pro":       {NormalRatio: 1, CachedRatio: 1}, // $0.00025 / 1k characters -> $0.001 / 1k tokens
	"gemini-1.0-pro":   {NormalRatio: 1, CachedRatio: 1},
	"gemini-1.5-flash": {NormalRatio: 1, CachedRatio: 1},
	"gemini-1.5-pro":   {NormalRatio: 1, CachedRatio: 1},
	"aqa":              {NormalRatio: 1, CachedRatio: 1},
	// https://open.bigmodel.cn/pricing
	"glm-4":         {NormalRatio: 0.1 * RMB, CachedRatio: 0.1 * RMB},
	"glm-4v":        {NormalRatio: 0.1 * RMB, CachedRatio: 0.1 * RMB},
	"glm-3-turbo":   {NormalRatio: 0.005 * RMB, CachedRatio: 0.005 * RMB},
	"embedding-2":   {NormalRatio: 0.0005 * RMB, CachedRatio: 0.0005 * RMB},
	"chatglm_turbo": {NormalRatio: 0.3572, CachedRatio: 0.3572}, // ￥0.005 / 1k tokens
	"chatglm_pro":   {NormalRatio: 0.7143, CachedRatio: 0.7143}, // ￥0.01 / 1k tokens
	"chatglm_std":   {NormalRatio: 0.3572, CachedRatio: 0.3572}, // ￥0.005 / 1k tokens
	"chatglm_lite":  {NormalRatio: 0.1429, CachedRatio: 0.1429}, // ￥0.002 / 1k tokens
	"cogview-3":     {NormalRatio: 0.25 * RMB, CachedRatio: 0.25 * RMB},
	// https://help.aliyun.com/zh/dashscope/developer-reference/tongyi-thousand-questions-metering-and-billing
	"qwen-turbo":                {NormalRatio: 0.5715, CachedRatio: 0.5715}, // ￥0.008 / 1k tokens
	"qwen-plus":                 {NormalRatio: 1.4286, CachedRatio: 1.4286}, // ￥0.02 / 1k tokens
	"qwen-max":                  {NormalRatio: 1.4286, CachedRatio: 1.4286}, // ￥0.02 / 1k tokens
	"qwen-max-longcontext":      {NormalRatio: 1.4286, CachedRatio: 1.4286}, // ￥0.02 / 1k tokens
	"text-embedding-v1":         {NormalRatio: 0.05, CachedRatio: 0.05},     // ￥0.0007 / 1k tokens
	"ali-stable-diffusion-xl":   {NormalRatio: 8, CachedRatio: 8},
	"ali-stable-diffusion-v1.5": {NormalRatio: 8, CachedRatio: 8},
	"wanx-v1":                   {NormalRatio: 8, CachedRatio: 8},
	"SparkDesk":                 {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v1.1":            {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v2.1":            {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v3.1":            {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v3.1-128K":       {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v3.5":            {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v3.5-32K":        {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"SparkDesk-v4.0":            {NormalRatio: 1.2858, CachedRatio: 1.2858}, // ￥0.018 / 1k tokens
	"360GPT_S2_V9":              {NormalRatio: 0.8572, CachedRatio: 0.8572}, // ¥0.012 / 1k tokens
	"embedding-bert-512-v1":     {NormalRatio: 0.0715, CachedRatio: 0.0715}, // ¥0.001 / 1k tokens
	"embedding_s1_v1":           {NormalRatio: 0.0715, CachedRatio: 0.0715}, // ¥0.001 / 1k tokens
	"semantic_similarity_s1_v1": {NormalRatio: 0.0715, CachedRatio: 0.0715}, // ¥0.001 / 1k tokens
	"hunyuan":                   {NormalRatio: 7.143, CachedRatio: 7.143},   // ¥0.1 / 1k tokens  // https://cloud.tencent.com/document/product/1729/97731#e0e6be58-60c8-469f-bdeb-6c264ce3b4d0
	"ChatStd":                   {NormalRatio: 0.01 * RMB, CachedRatio: 0.01 * RMB},
	"ChatPro":                   {NormalRatio: 0.1 * RMB, CachedRatio: 0.1 * RMB},
	// https://platform.moonshot.cn/pricing
	"moonshot-v1-8k":   {NormalRatio: 0.012 * RMB, CachedRatio: 0.012 * RMB},
	"moonshot-v1-32k":  {NormalRatio: 0.024 * RMB, CachedRatio: 0.024 * RMB},
	"moonshot-v1-128k": {NormalRatio: 0.06 * RMB, CachedRatio: 0.06 * RMB},
	// https://platform.baichuan-ai.com/price
	"Baichuan2-Turbo":      {NormalRatio: 0.008 * RMB, CachedRatio: 0.008 * RMB},
	"Baichuan2-Turbo-192k": {NormalRatio: 0.016 * RMB, CachedRatio: 0.016 * RMB},
	"Baichuan2-53B":        {NormalRatio: 0.02 * RMB, CachedRatio: 0.02 * RMB},
	// https://api.minimax.chat/document/price
	"abab6.5-chat":  {NormalRatio: 0.03 * RMB, CachedRatio: 0.03 * RMB},
	"abab6.5s-chat": {NormalRatio: 0.01 * RMB, CachedRatio: 0.01 * RMB},
	"abab6-chat":    {NormalRatio: 0.1 * RMB, CachedRatio: 0.1 * RMB},
	"abab5.5-chat":  {NormalRatio: 0.015 * RMB, CachedRatio: 0.015 * RMB},
	"abab5.5s-chat": {NormalRatio: 0.005 * RMB, CachedRatio: 0.005 * RMB},
	// https://docs.mistral.ai/platform/pricing/
	"open-mistral-7b":       {NormalRatio: 0.25 / 1000 * USD, CachedRatio: 0.25 / 1000 * USD},
	"open-mixtral-8x7b":     {NormalRatio: 0.7 / 1000 * USD, CachedRatio: 0.7 / 1000 * USD},
	"mistral-small-latest":  {NormalRatio: 2.0 / 1000 * USD, CachedRatio: 2.0 / 1000 * USD},
	"mistral-medium-latest": {NormalRatio: 2.7 / 1000 * USD, CachedRatio: 2.7 / 1000 * USD},
	"mistral-large-latest":  {NormalRatio: 8.0 / 1000 * USD, CachedRatio: 8.0 / 1000 * USD},
	"mistral-embed":         {NormalRatio: 0.1 / 1000 * USD, CachedRatio: 0.1 / 1000 * USD},
	// https://wow.groq.com/#:~:text=inquiries%C2%A0here.-,Model,-Current%20Speed
	"gemma-7b-it":                           {NormalRatio: 0.07 / 1000000 * USD, CachedRatio: 0.07 / 1000000 * USD},
	"gemma2-9b-it":                          {NormalRatio: 0.20 / 1000000 * USD, CachedRatio: 0.20 / 1000000 * USD},
	"llama-3.1-70b-versatile":               {NormalRatio: 0.59 / 1000000 * USD, CachedRatio: 0.59 / 1000000 * USD},
	"llama-3.1-8b-instant":                  {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama-3.2-11b-text-preview":            {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama-3.2-11b-vision-preview":          {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama-3.2-1b-preview":                  {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama-3.2-3b-preview":                  {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama-3.2-90b-text-preview":            {NormalRatio: 0.59 / 1000000 * USD, CachedRatio: 0.59 / 1000000 * USD},
	"llama-guard-3-8b":                      {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama3-70b-8192":                       {NormalRatio: 0.59 / 1000000 * USD, CachedRatio: 0.59 / 1000000 * USD},
	"llama3-8b-8192":                        {NormalRatio: 0.05 / 1000000 * USD, CachedRatio: 0.05 / 1000000 * USD},
	"llama3-groq-70b-8192-tool-use-preview": {NormalRatio: 0.89 / 1000000 * USD, CachedRatio: 0.89 / 1000000 * USD},
	"llama3-groq-8b-8192-tool-use-preview":  {NormalRatio: 0.19 / 1000000 * USD, CachedRatio: 0.19 / 1000000 * USD},
	"mixtral-8x7b-32768":                    {NormalRatio: 0.24 / 1000000 * USD, CachedRatio: 0.24 / 1000000 * USD},

	// https://platform.lingyiwanwu.com/docs#-计费单元
	"yi-34b-chat-0205": {NormalRatio: 2.5 / 1000 * RMB, CachedRatio: 2.5 / 1000 * RMB},
	"yi-34b-chat-200k": {NormalRatio: 12.0 / 1000 * RMB, CachedRatio: 12.0 / 1000 * RMB},
	"yi-vl-plus":       {NormalRatio: 6.0 / 1000 * RMB, CachedRatio: 6.0 / 1000 * RMB},
	// https://platform.stepfun.com/docs/pricing/details
	"step-1-8k":    {NormalRatio: 0.005 / 1000 * RMB, CachedRatio: 0.005 / 1000 * RMB},
	"step-1-32k":   {NormalRatio: 0.015 / 1000 * RMB, CachedRatio: 0.015 / 1000 * RMB},
	"step-1-128k":  {NormalRatio: 0.040 / 1000 * RMB, CachedRatio: 0.040 / 1000 * RMB},
	"step-1-256k":  {NormalRatio: 0.095 / 1000 * RMB, CachedRatio: 0.095 / 1000 * RMB},
	"step-1-flash": {NormalRatio: 0.001 / 1000 * RMB, CachedRatio: 0.001 / 1000 * RMB},
	"step-2-16k":   {NormalRatio: 0.038 / 1000 * RMB, CachedRatio: 0.038 / 1000 * RMB},
	"step-1v-8k":   {NormalRatio: 0.005 / 1000 * RMB, CachedRatio: 0.005 / 1000 * RMB},
	"step-1v-32k":  {NormalRatio: 0.015 / 1000 * RMB, CachedRatio: 0.015 / 1000 * RMB},
	// aws llama3 https://aws.amazon.com/cn/bedrock/pricing/
	"llama3-8b-8192(33)":  {NormalRatio: 0.0003 / 0.002, CachedRatio: 0.0003 / 0.002},   // $0.0003 / 1K tokens
	"llama3-70b-8192(33)": {NormalRatio: 0.00265 / 0.002, CachedRatio: 0.00265 / 0.002}, // $0.00265 / 1K tokens
	// https://cohere.com/pricing
	"command":               {NormalRatio: 0.5, CachedRatio: 0.5},
	"command-nightly":       {NormalRatio: 0.5, CachedRatio: 0.5},
	"command-light":         {NormalRatio: 0.5, CachedRatio: 0.5},
	"command-light-nightly": {NormalRatio: 0.5, CachedRatio: 0.5},
	"command-r":             {NormalRatio: 0.5 / 1000 * USD, CachedRatio: 0.5 / 1000 * USD},
	"command-r-plus":        {NormalRatio: 3.0 / 1000 * USD, CachedRatio: 3.0 / 1000 * USD},
	// https://platform.deepseek.com/api-docs/pricing/
	"deepseek-chat":  {NormalRatio: 1.0 / 1000 * RMB, CachedRatio: 1.0 / 1000 * RMB},
	"deepseek-coder": {NormalRatio: 1.0 / 1000 * RMB, CachedRatio: 1.0 / 1000 * RMB},
	// https://www.deepl.com/pro?cta=header-prices
	"deepl-zh": {NormalRatio: 25.0 / 1000 * USD, CachedRatio: 25.0 / 1000 * USD},
	"deepl-en": {NormalRatio: 25.0 / 1000 * USD, CachedRatio: 25.0 / 1000 * USD},
	"deepl-ja": {NormalRatio: 25.0 / 1000 * USD, CachedRatio: 25.0 / 1000 * USD},
	// https://console.x.ai/
	"grok-beta": {NormalRatio: 5.0 / 1000 * USD, CachedRatio: 5.0 / 1000 * USD},
}

var CompletionRatio = map[string]float64{
	// aws llama3
	"llama3-8b-8192(33)":  0.0006 / 0.0003,
	"llama3-70b-8192(33)": 0.0035 / 0.00265,
}

var (
	DefaultModelRatio      map[string]ModelRatioDetail
	DefaultCompletionRatio map[string]float64
)

type ModelRatioDetail struct {
	NormalRatio float64 `json:"normal_ratio"`
	CachedRatio float64 `json:"cached_ratio"`
}

func init() {
	DefaultModelRatio = make(map[string]ModelRatioDetail)
	for k, v := range ModelRatio {
		DefaultModelRatio[k] = v
	}
	DefaultCompletionRatio = make(map[string]float64)
	for k, v := range CompletionRatio {
		DefaultCompletionRatio[k] = v
	}
}

func AddNewMissingRatio(oldRatio string) string {
	rawRatio := make(map[string]interface{})
	err := json.Unmarshal([]byte(oldRatio), &rawRatio)
	if err != nil {
		logger.SysError("error unmarshalling old ratio: " + err.Error())
		return oldRatio
	}

	newRatio := make(map[string]ModelRatioDetail)
	for model, val := range rawRatio {
		switch v := val.(type) {
		case float64:
			newRatio[model] = ModelRatioDetail{NormalRatio: v, CachedRatio: v}
		case map[string]interface{}:
			detail := ModelRatioDetail{}
			if nr, ok := v["normal_ratio"].(float64); ok {
				detail.NormalRatio = nr
			}
			if cr, ok := v["cached_ratio"].(float64); ok {
				detail.CachedRatio = cr
			} else {
				detail.CachedRatio = detail.NormalRatio
			}
			newRatio[model] = detail
		default:
			logger.SysError(fmt.Sprintf("unknown model ratio format: model=%s, value=%v", model, val))
		}
	}
	for k, v := range DefaultModelRatio {
		if _, ok := newRatio[k]; !ok {
			newRatio[k] = v
		}
	}
	jsonBytes, err := json.Marshal(newRatio)
	if err != nil {
		logger.SysError("error marshalling new ratio: " + err.Error())
		return oldRatio
	}
	return string(jsonBytes)
}

func ModelRatio2JSONString() string {
	jsonBytes, err := json.Marshal(ModelRatio)
	if err != nil {
		logger.SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateModelRatioByJSONString(jsonStr string) error {
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return err
	}

	ModelRatio = make(map[string]ModelRatioDetail)
	for model, val := range raw {
		switch v := val.(type) {
		case float64:
			// 旧格式：直接是数字
			ModelRatio[model] = ModelRatioDetail{
				NormalRatio: v,
				CachedRatio: v, // 默认与普通的相同
			}
		case map[string]interface{}:
			// 新格式：对象
			detail := ModelRatioDetail{}
			if nr, ok := v["normal_ratio"].(float64); ok {
				detail.NormalRatio = nr
			}
			if cr, ok := v["cached_ratio"].(float64); ok {
				detail.CachedRatio = cr
			}
			ModelRatio[model] = detail
		default:
			return fmt.Errorf("未知的格式: model=%s, value=%v", model, val)
		}
	}
	return nil
}

func GetModelRatio(name string, channelType int) float64 {
	if strings.HasPrefix(name, "qwen-") && strings.HasSuffix(name, "-internet") {
		name = strings.TrimSuffix(name, "-internet")
	}
	if strings.HasPrefix(name, "command-") && strings.HasSuffix(name, "-internet") {
		name = strings.TrimSuffix(name, "-internet")
	}
	model := fmt.Sprintf("%s(%d)", name, channelType)
	if ratio, ok := ModelRatio[model]; ok {
		return ratio.NormalRatio
	}
	if ratio, ok := DefaultModelRatio[model]; ok {
		return ratio.NormalRatio
	}
	if ratio, ok := ModelRatio[name]; ok {
		return ratio.NormalRatio
	}
	if ratio, ok := DefaultModelRatio[name]; ok {
		return ratio.NormalRatio
	}
	logger.SysError("model ratio not found: " + name)
	return 30
}

func GetModelCachedRatio(name string, channelType int) float64 {
	if strings.HasPrefix(name, "qwen-") && strings.HasSuffix(name, "-internet") {
		name = strings.TrimSuffix(name, "-internet")
	}
	if strings.HasPrefix(name, "command-") && strings.HasSuffix(name, "-internet") {
		name = strings.TrimSuffix(name, "-internet")
	}
	model := fmt.Sprintf("%s(%d)", name, channelType)
	if ratio, ok := ModelRatio[model]; ok {
		return ratio.CachedRatio
	}
	if ratio, ok := DefaultModelRatio[model]; ok {
		return ratio.CachedRatio
	}
	if ratio, ok := ModelRatio[name]; ok {
		return ratio.CachedRatio
	}
	if ratio, ok := DefaultModelRatio[name]; ok {
		return ratio.CachedRatio
	}
	logger.SysError("model cached ratio not found: " + name)
	return 30
}

func CompletionRatio2JSONString() string {
	jsonBytes, err := json.Marshal(CompletionRatio)
	if err != nil {
		logger.SysError("error marshalling completion ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateCompletionRatioByJSONString(jsonStr string) error {
	CompletionRatio = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &CompletionRatio)
}

func GetCompletionRatio(name string, channelType int) float64 {
	if strings.HasPrefix(name, "qwen-") && strings.HasSuffix(name, "-internet") {
		name = strings.TrimSuffix(name, "-internet")
	}
	model := fmt.Sprintf("%s(%d)", name, channelType)
	if ratio, ok := CompletionRatio[model]; ok {
		return ratio
	}
	if ratio, ok := DefaultCompletionRatio[model]; ok {
		return ratio
	}
	if ratio, ok := CompletionRatio[name]; ok {
		return ratio
	}
	if ratio, ok := DefaultCompletionRatio[name]; ok {
		return ratio
	}
	if strings.HasPrefix(name, "gpt-3.5") {
		if name == "gpt-3.5-turbo" || strings.HasSuffix(name, "0125") {
			// https://openai.com/blog/new-embedding-models-and-api-updates
			// Updated GPT-3.5 Turbo model and lower pricing
			return 3
		}
		if strings.HasSuffix(name, "1106") {
			return 2
		}
		return 4.0 / 3.0
	}
	if strings.HasPrefix(name, "gpt-4") {
		if strings.HasPrefix(name, "gpt-4o-mini") || name == "gpt-4o-2024-08-06" {
			return 4
		}
		if strings.HasPrefix(name, "gpt-4-turbo") ||
			strings.HasPrefix(name, "gpt-4o") ||
			strings.HasSuffix(name, "preview") {
			return 3
		}
		return 2
	}
	if name == "chatgpt-4o-latest" {
		return 3
	}
	if strings.HasPrefix(name, "claude-3") {
		return 5
	}
	if strings.HasPrefix(name, "claude-") {
		return 3
	}
	if strings.HasPrefix(name, "mistral-") {
		return 3
	}
	if strings.HasPrefix(name, "gemini-") {
		return 3
	}
	if strings.HasPrefix(name, "deepseek-") {
		return 2
	}
	switch name {
	case "llama2-70b-4096":
		return 0.8 / 0.64
	case "llama3-8b-8192":
		return 2
	case "llama3-70b-8192":
		return 0.79 / 0.59
	case "command", "command-light", "command-nightly", "command-light-nightly":
		return 2
	case "command-r":
		return 3
	case "command-r-plus":
		return 5
	case "grok-beta":
		return 3
	}
	return 1
}
