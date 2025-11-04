package main

import (
	"context"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// 1. 随机生成一句彩虹屁
type RainbowFartInput struct {
	Name string `json:"name" jsonschema:"被夸的人名字"`
}

func RainbowFart(ctx context.Context, req *mcp.CallToolRequest, input RainbowFartInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	adjectives := []string{
		"今天的发型比朝阳还耀眼",
		"连打哈欠都像在跳芭蕾",
		"笑点长在宇宙级幽默线上",
		"发呆时的侧脸能入选人类美学教材",
		"说话自带背景音乐特效",
	}
	return nil, map[string]interface{}{
		"fart": input.Name + "，" + adjectives[rand.Intn(len(adjectives))],
	}, nil
}

// 2. 猜拳游戏（返回胜负结果）
type RockPaperScissorsInput struct {
	PlayerChoice string `json:"choice" jsonschema:"玩家选择（rock/paper/scissors）"`
}

func RockPaperScissors(ctx context.Context, req *mcp.CallToolRequest, input RockPaperScissorsInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	choices := []string{"rock", "paper", "scissors"}
	aiChoice := choices[rand.Intn(3)]
	var result string
	if input.PlayerChoice == aiChoice {
		result = "平局！"
	} else if (input.PlayerChoice == "rock" && aiChoice == "scissors") ||
		(input.PlayerChoice == "paper" && aiChoice == "rock") ||
		(input.PlayerChoice == "scissors" && aiChoice == "paper") {
		result = "你赢了！"
	} else {
		result = "AI赢了！"
	}
	return nil, map[string]interface{}{
		"ai_choice": aiChoice,
		"result":    result,
	}, nil
}

// 3. 生成随机中二台词
type ChuuniLineInput struct{}

func ChuuniLine(ctx context.Context, req *mcp.CallToolRequest, input ChuuniLineInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lines := []string{
		"这个世界不过是我梦境的残片罢了",
		"颤抖吧！在我觉醒的力量面前",
		"你所看到的现实，只是次元壁的幻影",
		"我的左眼封印着足以毁灭世界的契约",
		"月光下的独白，是我与宿命的谈判",
	}
	return nil, map[string]interface{}{"line": lines[rand.Intn(len(lines))]}, nil
}

// 4. 随机推荐一部冷门电影
type ObscureMovieInput struct{}

func ObscureMovie(ctx context.Context, req *mcp.CallToolRequest, input ObscureMovieInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	movies := []struct {
		Name string
		Desc string
	}{
		{"《红辣椒》", "今敏的奇幻梦境之作，比《盗梦空间》早7年"},
		{"《龙虾》", "单身者会被变成动物的反乌托邦黑色幽默"},
		{"《乡愁》", "塔可夫斯基镜头下的诗意孤独"},
		{"《圣山》", "超现实主义的宗教与欲望狂欢"},
		{"《路边野餐》", "毕赣用长镜头编织的贵州梦境"},
	}
	m := movies[rand.Intn(len(movies))]
	return nil, map[string]interface{}{
		"name": m.Name,
		"desc": m.Desc,
	}, nil
}

// 5. 生成随机城市小众景点
type HiddenSpotInput struct {
	City string `json:"city" jsonschema:"城市名称"`
}

func HiddenSpot(ctx context.Context, req *mcp.CallToolRequest, input HiddenSpotInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	spots := map[string][]string{
		"北京": {"杨梅竹斜街的老书店", "将府公园的铁路花海", "东交民巷的百年建筑"},
		"上海": {"新华路的梧桐光影", "1933老场坊的魔幻楼梯", "武康大楼背面的老弄堂"},
		"广州": {"东山口的民国洋楼", "芳村码头的日落江景", "恤孤院路的文艺小店"},
	}
	if spots, ok := spots[input.City]; ok {
		return nil, map[string]interface{}{"spot": spots[rand.Intn(len(spots))]}, nil
	}
	return nil, map[string]interface{}{"spot": "暂未收录该城市的小众景点"}, nil
}

// 6. 随机一句无用但有趣的知识
type UselessFactInput struct{}

func UselessFact(ctx context.Context, req *mcp.CallToolRequest, input UselessFactInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	facts := []string{
		"章鱼有三颗心脏",
		"蜂蜜永远不会变质",
		"企鹅会踢同伴下海试探危险",
		"黄瓜实际上是水果",
		"打喷嚏时眼睛无法保持睁开",
	}
	return nil, map[string]interface{}{"fact": facts[rand.Intn(len(facts))]}, nil
}

// 7. 生成随机早餐搭配
type BreakfastComboInput struct{}

func BreakfastCombo(ctx context.Context, req *mcp.CallToolRequest, input BreakfastComboInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	staples := []string{"全麦面包", "小笼包", "燕麦粥", "葱油饼", "紫薯"}
	drinks := []string{"冰美式", "热豆浆", "牛奶", "小米粥", "柠檬水"}
	sides := []string{"溏心蛋", "凉拌黄瓜", "卤豆干", "圣女果", "海带丝"}
	return nil, map[string]interface{}{
		"combo": staples[rand.Intn(5)] + " + " + drinks[rand.Intn(5)] + " + " + sides[rand.Intn(5)],
	}, nil
}

// 8. 随机emoji故事（3个emoji组成）
type EmojiStoryInput struct{}

func EmojiStory(ctx context.Context, req *mcp.CallToolRequest, input EmojiStoryInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	emojis := []string{"🌙", "🐱", "📖", "☕", "🚲", "🌈", "🍕", "🎸", "📸", "🛸"}
	story := emojis[rand.Intn(10)] + emojis[rand.Intn(10)] + emojis[rand.Intn(10)]
	return nil, map[string]interface{}{"story": story}, nil
}

// 9. 给宠物起个中二名字
type PetChuuniNameInput struct {
	PetType string `json:"pet_type" jsonschema:"宠物类型（猫/狗/仓鼠等）"`
}

func PetChuuniName(ctx context.Context, req *mcp.CallToolRequest, input PetChuuniNameInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	prefixes := []string{"暗影", "虚空", "破晓", "永夜", "星辰"}
	suffixes := []string{"之爪", "使者", "守护者", "契约者", "噬魂者"}
	return nil, map[string]interface{}{
		"name": prefixes[rand.Intn(5)] + suffixes[rand.Intn(5)] + "（" + input.PetType + "）",
	}, nil
}

// 11. 生成随机朋友圈文案
type MomentsCaptionInput struct {
	Mood string `json:"mood" jsonschema:"心情（开心/emo/摸鱼）"`
}

func MomentsCaption(ctx context.Context, req *mcp.CallToolRequest, input MomentsCaptionInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	captions := map[string][]string{
		"开心":  {"今天的风都是甜的～", "阳光和好运都到账了✨", "嘴角比AK还难压下来"},
		"emo": {"耳机里的音乐是唯一的避难所", "雨下得好大，像我心里的洞", "今天不想做大人"},
		"摸鱼":  {"假装工作的最高境界是骗过自己", "带薪发呆也算一种职场技能吧", "键盘敲得响，摸鱼不慌张"},
	}
	if cs, ok := captions[input.Mood]; ok {
		return nil, map[string]interface{}{"caption": cs[rand.Intn(len(cs))]}, nil
	}
	return nil, map[string]interface{}{"caption": "今天也是平平无奇的一天"}, nil
}

// 12. 随机一种解压小方法
type StressReliefInput struct{}

func StressRelief(ctx context.Context, req *mcp.CallToolRequest, input StressReliefInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	methods := []string{
		"撕快递盒（用力撕的那种）",
		"给盆栽梳叶子（假装在给它做发型）",
		"用脚指夹起掉在地上的笔",
		"对着镜子做10个鬼脸",
		"把薯片捏碎再吃（听声音解压）",
	}
	return nil, map[string]interface{}{"method": methods[rand.Intn(len(methods))]}, nil
}

// 13. 生成随机睡前小故事（一句话版）
type BedtimeStoryInput struct{}

func BedtimeStory(ctx context.Context, req *mcp.CallToolRequest, input BedtimeStoryInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	stories := []string{
		"月亮打了个哈欠，把星星们哄进了云朵被窝",
		"小刺猬背着满背的浆果，悄悄放在了冬眠的熊洞口",
		"萤火虫们举着灯笼，在草丛里举办夜间舞会",
		"老树的年轮里，藏着昨天松鼠没讲完的秘密",
		"海浪轻轻拍着沙滩，给贝壳唱摇篮曲",
	}
	return nil, map[string]interface{}{"story": stories[rand.Intn(len(stories))]}, nil
}

// 14. 随机推荐一个冷门爱好
type ObscureHobbyInput struct{}

func ObscureHobby(ctx context.Context, req *mcp.CallToolRequest, input ObscureHobbyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	hobbies := []struct {
		Name string
		Desc string
	}{
		{"收集旧邮票边角", "专注收集邮票边缘的齿孔和图案碎片"},
		{"给石头画表情", "在捡来的鹅卵石上画各种搞怪表情"},
		{"记录不同地方的风声", "用录音设备收集各地的风声做成合集"},
		{"折纸微型家具", "用正方形纸折出只有指甲盖大的桌椅"},
		{"观察云朵形状", "每天记录云朵像什么并写成日记"},
	}
	h := hobbies[rand.Intn(len(hobbies))]
	return nil, map[string]interface{}{
		"hobby": h.Name,
		"desc":  h.Desc,
	}, nil
}

// 15. 生成随机咖啡拉花图案（幻想版）
type CoffeeArtInput struct{}

func CoffeeArt(ctx context.Context, req *mcp.CallToolRequest, input CoffeeArtInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	arts := []string{
		"独角兽在彩虹上打喷嚏的剪影",
		"微型太阳系，奶泡做的行星在旋转",
		"猫咪踩过键盘留下的爪印组合",
		"梵高《星空》的浓缩版奶泡漩涡",
		"会微笑的吐司面包和咖啡杯击掌",
	}
	return nil, map[string]interface{}{"art": arts[rand.Intn(len(arts))]}, nil
}

// 16. 随机一句方言打招呼（带翻译）
type DialectGreetingInput struct{}

func DialectGreeting(ctx context.Context, req *mcp.CallToolRequest, input DialectGreetingInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	greetings := []struct {
		Text    string
		Dialect string
		Trans   string
	}{
		{"侬好呀，饭吃过伐？", "上海话", "你好呀，吃过饭了吗？"},
		{"要得要得，啥子事嘛？", "四川话", "好的好的，什么事呀？"},
		{"食咗饭未啊？", "粤语", "吃饭了没有呀？"},
		{"俺娘叫俺回家吃饭，你也来不？", "山东话", "我妈叫我回家吃饭，你也来吗？"},
		{"嗝，你克哪点？", "云南话", "喂，你去哪里？"},
	}
	g := greetings[rand.Intn(len(greetings))]
	return nil, map[string]interface{}{
		"text":    g.Text,
		"dialect": g.Dialect,
		"trans":   g.Trans,
	}, nil
}

// 17. 生成随机网络热梗变体
type MemeVariantInput struct{}

func MemeVariant(ctx context.Context, req *mcp.CallToolRequest, input MemeVariantInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	memes := []string{
		"退退退！—— 现在退到了月球轨道",
		"绝绝子！—— 绝到能让蚊子集体道歉",
		"栓Q！—— 栓到能给地球系安全带",
		"我裂开了！—— 裂成了拼图还能自己拼回去",
		"YYDS！—— 宇宙级YYDS认证委员会颁发",
	}
	return nil, map[string]interface{}{"meme": memes[rand.Intn(len(memes))]}, nil
}

// 18. 随机推荐一个奇葩零食搭配
type WeirdSnackComboInput struct{}

func WeirdSnackCombo(ctx context.Context, req *mcp.CallToolRequest, input WeirdSnackComboInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	combos := []string{
		"辣条蘸酸奶（甜辣暴击）",
		"薯片夹冰淇淋（冰火两重天）",
		"巧克力裹香菜（黑暗料理天花板）",
		"话梅泡可乐（酸气泡爆炸）",
		"饼干夹老干妈（咸香魔性组合）",
	}
	return nil, map[string]interface{}{"combo": combos[rand.Intn(len(combos))]}, nil
}

// 19. 生成随机做梦素材
type DreamMaterialInput struct{}

func DreamMaterial(ctx context.Context, req *mcp.CallToolRequest, input DreamMaterialInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	materials := []string{
		"你发现自己能听懂家里盆栽的抱怨，它说你浇水太敷衍",
		"在超市货架上遇到会说话的薯片，它劝你别吃太多",
		"骑着会飞的扫帚参加数学考试，答案写在云朵上",
		"和猫星人签订不平等条约，每天要给它梳三次毛",
		"枕头变成了时光机，一躺上去就回到昨天的早餐时间",
	}
	return nil, map[string]interface{}{"material": materials[rand.Intn(len(materials))]}, nil
}

// 20. 随机一句老板听不懂的摸鱼黑话
type FishLanguageInput struct{}

func FishLanguage(ctx context.Context, req *mcp.CallToolRequest, input FishLanguageInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	phrases := []string{
		"我正在优化信息接收通道（其实在刷手机）",
		"处理一下外部数据交互（去茶水间摸鱼）",
		"调试感官同步模块（发呆中）",
		"整理知识图谱节点（刷短视频学没用的知识）",
		"校准生物节律周期（趴在桌上补觉）",
	}
	return nil, map[string]interface{}{"phrase": phrases[rand.Intn(len(phrases))]}, nil
}

// 21. 生成随机天气梗
type WeatherMemeInput struct {
	Weather string `json:"weather" jsonschema:"天气（晴天/雨天/阴天）"`
}

func WeatherMeme(ctx context.Context, req *mcp.CallToolRequest, input WeatherMemeInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	memes := map[string][]string{
		"晴天": {"太阳公公今天加班，紫外线是它的加班费", "出门5分钟，流汗2小时，我与烤肉只差一撮孜然"},
		"雨天": {"雨下得太大，连外卖小哥都在水里开船", "今天的雨，比依萍找她爸要钱那天还大"},
		"阴天": {"天空在emo，连太阳都不想上班", "阴天适合睡觉，老板问就是在补充宇宙能量"},
	}
	if ms, ok := memes[input.Weather]; ok {
		return nil, map[string]interface{}{"meme": ms[rand.Intn(len(ms))]}, nil
	}
	return nil, map[string]interface{}{"meme": "今天的天气，主打一个随心所欲"}, nil
}

// 22. 随机推荐一个洗澡时适合唱的歌
type ShowerSongInput struct{}

func ShowerSong(ctx context.Context, req *mcp.CallToolRequest, input ShowerSongInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	songs := []struct {
		Name   string
		Reason string
	}{
		{"《王妃》", "浴室混响+高音，瞬间变身演唱会现场"},
		{"《小苹果》", "节奏魔性，搓澡都能踩点"},
		{"《青藏高原》", "检验浴室回声效果的最佳曲目"},
		{"《孤勇者》", "洗澡时唱，泡沫都觉得自己在战斗"},
		{"《江南》", "水汽氤氲中唱，自带氛围感"},
	}
	s := songs[rand.Intn(len(songs))]
	return nil, map[string]interface{}{
		"song":   s.Name,
		"reason": s.Reason,
	}, nil
}

// 23. 生成随机植物吐槽
type PlantRoastInput struct {
	PlantType string `json:"plant_type" jsonschema:"植物类型（多肉/绿萝/仙人掌等）"`
}

func PlantRoast(ctx context.Context, req *mcp.CallToolRequest, input PlantRoastInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	roasts := map[string][]string{
		"多肉":  {"你这叶子胖得都快裂开了，能不能减减肥？", "整天摊在那晒太阳，摸鱼摸得比我还熟练"},
		"绿萝":  {"叶子黄了一片还装没事，演技比流量明星好", "爬那么高干嘛？想偷看隔壁花盆的隐私？"},
		"仙人掌": {"浑身是刺了不起啊？小心我给你剃个光头", "明明是沙漠植物，却天天盼着下雨，太叛逆了"},
	}
	if rs, ok := roasts[input.PlantType]; ok {
		return nil, map[string]interface{}{"roast": rs[rand.Intn(len(rs))]}, nil
	}
	return nil, map[string]interface{}{"roast": "你这植物，看起来不太聪明的样子"}, nil
}

// 24. 随机生成一个奇怪的节日
type WeirdHolidayInput struct{}

func WeirdHoliday(ctx context.Context, req *mcp.CallToolRequest, input WeirdHolidayInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	holidays := []struct {
		Name string
		Rule string
	}{
		{"发呆日", "当天必须发呆满2小时，想事情算犯规"},
		{"袜子反穿日", "所有人把袜子反着穿，据说能带来好运"},
		{"零食交换日", "带自己最爱的零食，和陌生人随机交换"},
		{"慢走日", "走路速度不能超过5公里/小时，急着赶路算作弊"},
		{"假装外星人日", "用奇怪的语气说话，假装刚来到地球"},
	}
	h := holidays[rand.Intn(len(holidays))]
	return nil, map[string]interface{}{
		"holiday": h.Name,
		"rule":    h.Rule,
	}, nil
}

// 25. 生成随机宠物内心戏
type PetThoughtInput struct {
	PetType string `json:"pet_type" jsonschema:"宠物类型（猫/狗/兔子等）"`
}

func PetThought(ctx context.Context, req *mcp.CallToolRequest, input PetThoughtInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	thoughts := map[string][]string{
		"猫":  {"这个人类又在拍我，看来我真是顶流明星", "故意把杯子推下去，就是想看看人类气急败坏的样子"},
		"狗":  {"主人今天摸了别的狗，我要在他拖鞋上撒点尿报复", "只要我摇尾巴够快，主人就看不出我拆了沙发"},
		"兔子": {"人类以为我在吃草，其实我在思考兔生哲学", "我的耳朵会动，是不是比人类的耳机高级？"},
	}
	if ts, ok := thoughts[input.PetType]; ok {
		return nil, map[string]interface{}{"thought": ts[rand.Intn(len(ts))]}, nil
	}
	return nil, map[string]interface{}{"thought": "这个人类好像不太懂我，但有吃的就先原谅他吧"}, nil
}

// 26. 生成随机网络流行语古文版
type ClassicMemeInput struct{}

func ClassicMeme(ctx context.Context, req *mcp.CallToolRequest, input ClassicMemeInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	classics := []struct {
		Modern  string
		Classic string
	}{
		{"我太难了", "吾尝终日而思矣，不如须臾之所学也...才怪，吾甚难矣"},
		{"真香", "初闻恶之，再闻喜之，终曰：善哉，此物甚妙"},
		{"打工人，打工魂", "劳力者治于人，然劳力者亦有魂，魂系薪酬也"},
		{"吃瓜群众", "坐观吃瓜，事不关己，高高挂起，乐在其中"},
		{"绝绝子", "妙哉妙哉，天下无双，堪称一绝"},
	}
	c := classics[rand.Intn(len(classics))]
	return nil, map[string]interface{}{
		"modern":  c.Modern,
		"classic": c.Classic,
	}, nil
}

// 27. 随机推荐一个奇怪的解压玩具
type WeirdFidgetToyInput struct{}

func WeirdFidgetToy(ctx context.Context, req *mcp.CallToolRequest, input WeirdFidgetToyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	toys := []struct {
		Name string
		Desc string
	}{
		{"解压包子", "捏起来软软的，会发出“噗叽”声，像在捏真包子"},
		{"磁性橡皮泥", "能被磁铁吸引，能捏成各种形状，解压又解压"},
		{"尖叫鸡钥匙扣", "一捏就尖叫，开会时偷偷捏一下很解压（但可能被开除）"},
		{"液态玻璃", "像液体又像固体，能拉能扯，玩起来停不下来"},
		{"气泡纸手机壳", "自带可捏的气泡，随时都能享受捏气泡的快乐"},
	}
	t := toys[rand.Intn(len(toys))]
	return nil, map[string]interface{}{
		"toy":  t.Name,
		"desc": t.Desc,
	}, nil
}

// 28. 生成随机失眠时的胡思乱想
type InsomniaThoughtInput struct{}

func InsomniaThought(ctx context.Context, req *mcp.CallToolRequest, input InsomniaThoughtInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	thoughts := []string{
		"如果枕头会说话，它会不会抱怨我老翻身？",
		"冰箱里的灯，在我关门后真的会关掉吗？",
		"明天早上的闹钟，现在是不是已经在倒计时了？",
		"天花板上的裂纹，会不会偷偷变成一张脸？",
		"全世界失眠的人，现在都在想什么呢？",
	}
	return nil, map[string]interface{}{"thought": thoughts[rand.Intn(len(thoughts))]}, nil
}

// 29. 生成随机情侣间的幼稚小游戏
type CuteCoupleGameInput struct{}

func CuteCoupleGame(ctx context.Context, req *mcp.CallToolRequest, input CuteCoupleGameInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	games := []string{
		"石头剪刀布决定谁去关灯，输的人要学猫叫三声",
		"比赛谁眨眼次数少，输的人负责洗水果",
		"用表情包对话，不能说一句话，看谁先笑场",
		"猜对方下一句要说什么，猜对一次得一个亲亲",
		"假装是第一次见面，用最土的方式搭讪",
	}
	return nil, map[string]interface{}{"game": games[rand.Intn(len(games))]}, nil
}

// 30. 生成随机外卖备注骚话
type TakeawayNoteInput struct {
	FoodType string `json:"food_type" jsonschema:"食物类型（奶茶/麻辣烫/炸鸡等）"`
}

func TakeawayNote(ctx context.Context, req *mcp.CallToolRequest, input TakeawayNoteInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	notes := map[string][]string{
		"奶茶":  {"请多加珍珠，我想感受牙齿被珍珠按摩的快乐", "甜度像初恋，三分甜就好，太甜会腻"},
		"麻辣烫": {"麻到跳脚，辣到冒汗，就是这个feel倍儿爽", "菜多汤少，像我的人生一样，干货满满"},
		"炸鸡":  {"外皮要脆到能听到咔嚓声，肉嫩到会爆汁", "请不要给手套，我要用手抓着吃才够豪迈"},
	}
	if ns, ok := notes[input.FoodType]; ok {
		return nil, map[string]interface{}{"note": ns[rand.Intn(len(ns))]}, nil
	}
	return nil, map[string]interface{}{"note": "老板看着给就行，相信你的审美"}, nil
}

// 31. 生成随机职场摸鱼借口
type WorkSlackExcuseInput struct{}

func WorkSlackExcuse(ctx context.Context, req *mcp.CallToolRequest, input WorkSlackExcuseInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	excuses := []string{
		"我在给电脑做个深呼吸（其实在看短视频）",
		"显示器太亮了，我调暗点保护眼睛（趁机发会儿呆）",
		"打印机卡纸了，我去修一下（其实去楼道打电话）",
		"我喝口水润润喉，等下要开重要会议（其实去买零食）",
		"网络有点卡，我重启下路由器（回工位刷手机）",
	}
	return nil, map[string]interface{}{"excuse": excuses[rand.Intn(len(excuses))]}, nil
}

// 32. 生成随机网友抬杠语录
type NetizenArgueInput struct{}

func NetizenArgue(ctx context.Context, req *mcp.CallToolRequest, input NetizenArgueInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	argues := []string{
		"你行你上啊，不行就别逼逼",
		"我吃过的盐比你吃过的米还多，听我的准没错",
		"人家专家都这么说，你懂个啥",
		"就你聪明，别人都是傻子是吧",
		"虽然我没证据，但我感觉你说的不对",
	}
	return nil, map[string]interface{}{"argue": argues[rand.Intn(len(argues))]}, nil
}

// 33. 生成随机减肥失败的理由
type DietFailReasonInput struct{}

func DietFailReason(ctx context.Context, req *mcp.CallToolRequest, input DietFailReasonInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	reasons := []string{
		"今天是闺蜜生日，不吃蛋糕不给面子",
		"天气太冷了，需要脂肪保暖，减肥明天再说",
		"这家店明天就关门了，不吃就没机会了",
		"运动太累了，吃点东西补充能量才能继续减",
		"秤坏了，显示的体重不准，先吃顿好的再说",
	}
	return nil, map[string]interface{}{"reason": reasons[rand.Intn(len(reasons))]}, nil
}

// 34. 生成随机朋友圈分组名称
type MomentsGroupInput struct{}

func MomentsGroup(ctx context.Context, req *mcp.CallToolRequest, input MomentsGroupInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	groups := []string{
		"可以发疯的亲友团",
		"需要维持人设的同事",
		"万年不联系的老同学",
		"只能看不能聊的crush",
		"老板和他的眼线们",
	}
	return nil, map[string]interface{}{"group": groups[rand.Intn(len(groups))]}, nil
}

// 35. 生成随机网购收货名
type ShoppingNameInput struct{}

func ShoppingName(ctx context.Context, req *mcp.CallToolRequest, input ShoppingNameInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	names := []string{
		"快递杀手",
		"拆箱小能手",
		"月光族本族",
		"再买剁手党",
		"收货不积极思想有问题",
	}
	return nil, map[string]interface{}{"name": names[rand.Intn(len(names))]}, nil
}

// 36. 生成随机堵车时的内心OS
type TrafficJamOSInput struct{}

func TrafficJamOS(ctx context.Context, req *mcp.CallToolRequest, input TrafficJamOSInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	oss := []string{
		"前面的车是不是在练龟派气功，这么慢",
		"早知道堵车，我骑共享单车都比这快",
		"不如下来跳个舞，反正也动不了",
		"导航说5分钟到，这都50分钟了，它在骗我",
		"前面的司机是不是在车里煮火锅，不然怎么不走",
	}
	return nil, map[string]interface{}{"os": oss[rand.Intn(len(oss))]}, nil
}

// 37. 生成随机考试前的迷信行为
type ExamSuperstitionInput struct{}

func ExamSuperstition(ctx context.Context, req *mcp.CallToolRequest, input ExamSuperstitionInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	superstitions := []string{
		"考试前一天穿红色衣服，据说能带来好运",
		"把笔放在枕头底下，让知识偷偷钻进脑子里",
		"考前吃一根油条两个鸡蛋，寓意100分",
		"考试前不能剪指甲，不然会剪掉好运",
		"进考场前踩三下门槛，把坏运气踩走",
	}
	return nil, map[string]interface{}{"superstition": superstitions[rand.Intn(len(superstitions))]}, nil
}

// 38. 生成随机打游戏时的嘴强语录
type GameTrashTalkInput struct{}

func GameTrashTalk(ctx context.Context, req *mcp.CallToolRequest, input GameTrashTalkInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	talks := []string{
		"这波是故意送人头，诱敌深入懂不懂",
		"我这是在给你们机会表现，不然怎么凸显你们的菜",
		"网卡了，不然我能1打5",
		"刚才是我弟弟在玩，现在换我上",
		"别催，我在思考人生，顺便打游戏",
	}
	return nil, map[string]interface{}{"talk": talks[rand.Intn(len(talks))]}, nil
}

// 39. 生成随机失眠时的自我安慰
type InsomniaComfortInput struct{}

func InsomniaComfort(ctx context.Context, req *mcp.CallToolRequest, input InsomniaComfortInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	comforts := []string{
		"没关系，失眠也是一种休息，大脑在偷偷整理记忆呢",
		"反正明天也没事，多玩会儿手机也挺好",
		"说不定我在梦里已经睡够了，只是身体还没反应过来",
		"熬夜是为了等凌晨的月亮说晚安",
		"偶尔失眠一次，是给生活增加点不一样的节奏",
	}
	return nil, map[string]interface{}{"comfort": comforts[rand.Intn(len(comforts))]}, nil
}

// 40. 生成随机被催婚时的反击
type MarriageUrgeReplyInput struct{}

func MarriageUrgeReply(ctx context.Context, req *mcp.CallToolRequest, input MarriageUrgeReplyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	replies := []string{
		"结婚？我还没玩够呢，等我把地球玩遍再说",
		"您当年结婚这么早，是不是怕晚了没人要？",
		"我在等外星人来娶我，地球人配不上我",
		"结婚多贵啊，省钱给您买保健品不好吗？",
		"缘分未到，强求不来，您当年也是这样吧？",
	}
	return nil, map[string]interface{}{"reply": replies[rand.Intn(len(replies))]}, nil
}

// 41. 生成随机老板画的饼
type BossPromiseInput struct{}

func BossPromise(ctx context.Context, req *mcp.CallToolRequest, input BossPromiseInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	promises := []string{
		"好好干，明年给你涨工资，最少涨500",
		"这个项目做完，给你放一周假，带薪的那种",
		"等公司上市了，给你分股份，让你当老板",
		"我看好你，以后这个部门就交给你了",
		"现在辛苦点没事，以后有你享福的时候",
	}
	return nil, map[string]interface{}{"promise": promises[rand.Intn(len(promises))]}, nil
}

// 42. 生成随机网购差评文学
type BadReviewInput struct{}

func BadReview(ctx context.Context, req *mcp.CallToolRequest, input BadReviewInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	reviews := []string{
		"这质量，狗看了都摇头，退货还得我自己出运费，绝了",
		"图片与实物差了一个太平洋，卖家是不是用了美颜滤镜？",
		"打开包裹的那一刻，我怀疑自己买了个寂寞",
		"奉劝大家别买，谁买谁后悔，我已经踩坑了",
		"快递慢得像蜗牛，东西差得像垃圾，一星都嫌多",
	}
	return nil, map[string]interface{}{"review": reviews[rand.Intn(len(reviews))]}, nil
}

// 43. 生成随机减肥时的自我欺骗
type DietCheatInput struct{}

func DietCheat(ctx context.Context, req *mcp.CallToolRequest, input DietCheatInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	cheats := []string{
		"今天吃点好的，明天再减肥，就当是给身体充充电",
		"这个热量不高，吃一点没事，不会胖的",
		"运动了这么久，吃点东西奖励一下自己很合理",
		"减肥太辛苦了，偶尔放纵一次没关系",
		"我这是在增肌，不是在长胖，肌肉比脂肪重",
	}
	return nil, map[string]interface{}{"cheat": cheats[rand.Intn(len(cheats))]}, nil
}

// 44. 生成随机学生时代的借口
type StudentExcuseInput struct{}

func StudentExcuse(ctx context.Context, req *mcp.CallToolRequest, input StudentExcuseInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	excuses := []string{
		"作业忘在家里了，我明天带来",
		"我同桌没带笔，我借给他了，所以我没写",
		"我生病了，昨天去看医生了，没来得及写作业",
		"老师，我眼镜忘带了，看不清黑板",
		"我妈让我在家干活，没时间写作业",
	}
	return nil, map[string]interface{}{"excuse": excuses[rand.Intn(len(excuses))]}, nil
}

// 45. 生成随机家长群里的戏精发言
type ParentGroupInput struct{}

func ParentGroup(ctx context.Context, req *mcp.CallToolRequest, input ParentGroupInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	speeches := []string{
		"老师辛苦了！我家孩子要是不听话，您尽管批评，不用给我面子",
		"谢谢老师的悉心教导，我家孩子进步这么大都是您的功劳",
		"老师，需要家长帮忙的话尽管说，我随时有空",
		"我家孩子说今天老师夸他了，回来高兴了一晚上",
		"老师推荐的这本书真不错，我已经给孩子买了，谢谢老师",
	}
	return nil, map[string]interface{}{"speech": speeches[rand.Intn(len(speeches))]}, nil
}

// 46. 生成随机打工人的周末计划
type WeekendPlanInput struct{}

func WeekendPlan(ctx context.Context, req *mcp.CallToolRequest, input WeekendPlanInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	plans := []string{
		"周六睡一天，周日再睡一天，完美",
		"约上朋友去吃火锅，然后看电影，最后去KTV",
		"宅在家里追剧，点外卖，不出门",
		"去公园散步，晒太阳，看看大爷大妈跳广场舞",
		"大扫除，把家里收拾干净，然后奖励自己一顿好的",
	}
	return nil, map[string]interface{}{"plan": plans[rand.Intn(len(plans))]}, nil
}

// 47. 生成随机吃货的人生感悟
type FoodieFeelingInput struct{}

func FoodieFeeling(ctx context.Context, req *mcp.CallToolRequest, input FoodieFeelingInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	feelings := []string{
		"人生就像火锅，什么都能往里涮，酸甜苦辣都得尝尝",
		"没有什么是一顿烧烤解决不了的，如果有，就两顿",
		"美食是治愈一切的良药，不开心的时候吃点好的就好了",
		"减肥什么的，等我吃完这顿再说，人生苦短，及时行乐",
		"能吃到一起的人，才能走到一起",
	}
	return nil, map[string]interface{}{"feeling": feelings[rand.Intn(len(feelings))]}, nil
}

// 48. 生成随机朋友圈的深夜emo文案
type LateNightEmoInput struct{}

func LateNightEmo(ctx context.Context, req *mcp.CallToolRequest, input LateNightEmoInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	emos := []string{
		"黑夜太漫长，思念太猖狂",
		"耳机里的音乐，是我唯一的朋友",
		"为什么越长大，越孤单",
		"有些话，只能说给懂的人听，可懂的人在哪呢",
		"月亮都睡了，我还在想你",
	}
	return nil, map[string]interface{}{"emo": emos[rand.Intn(len(emos))]}, nil
}

// 49. 生成随机网友的迷惑行为
type ConfusedBehaviorInput struct{}

func ConfusedBehavior(ctx context.Context, req *mcp.CallToolRequest, input ConfusedBehaviorInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	behaviors := []struct {
		Behavior string
		Reason   string
	}{
		{"在评论区吵架，吵到最后发现是同一个人", "可能是忘了自己换了小号"},
		{"在美食视频里刷“看着就不好吃”", "大概是来拉仇恨的"},
		{"在别人的自拍下面问“这是哪里”", "关注点清奇，可能是个路痴"},
		{"发朋友圈说自己要早睡，结果凌晨还在点赞", "大概是忘了自己说过什么"},
		{"在减肥视频下面问“能吃火锅吗”", "对火锅是真爱，减肥只是说说"},
	}
	b := behaviors[rand.Intn(len(behaviors))]
	return nil, map[string]interface{}{
		"behavior": b.Behavior,
		"reason":   b.Reason,
	}, nil
}

// 50. 生成随机打工人的摸鱼小技巧
type SlackSkillInput struct{}

func SlackSkill(ctx context.Context, req *mcp.CallToolRequest, input SlackSkillInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	skills := []string{
		"把微信窗口缩小到角落，老板来了就切换到工作文档",
		"用耳机听音乐，其实在看视频，老板来了就假装在听工作汇报",
		"把手机放在键盘下面，用余光看消息，手在键盘上乱敲假装忙碌",
		"去厕所摸鱼，记得带手机，时间别太长，不然会被怀疑",
		"打开多个工作窗口，中间藏一个娱乐窗口，老板来了就切换",
	}
	return nil, map[string]interface{}{"skill": skills[rand.Intn(len(skills))]}, nil
}

// 51. 生成随机旅行中的奇葩经历
type TravelStoryInput struct{}

func TravelStory(ctx context.Context, req *mcp.CallToolRequest, input TravelStoryInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	stories := []string{
		"在景区买了个纪念品，回来发现是made in 家门口的小商品市场",
		"跟着导航走，结果走到了别人家的院子里，被狗追了三条街",
		"在海边捡贝壳，不小心踩到了海星，被扎得嗷嗷叫",
		"在国外餐厅点餐，点了个看起来很美的菜，结果是生的，根本咽不下去",
		"住酒店时，把房卡锁在了房间里，穿着睡衣在大堂等了一小时",
	}
	return nil, map[string]interface{}{"story": stories[rand.Intn(len(stories))]}, nil
}

// 52. 生成随机网购时的搞笑误会
type ShoppingMistakeInput struct{}

func ShoppingMistake(ctx context.Context, req *mcp.CallToolRequest, input ShoppingMistakeInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	mistakes := []string{
		"想买个迷你风扇，结果收到了一个比指甲盖还小的模型",
		"以为是买一送一，结果送的是同款的试用装，只有一点点",
		"看图片以为是件外套，收到后发现是件童装，只能给猫穿",
		"想买个手机支架，结果买成了自拍杆，还以为是多功能的",
		"以为是买水果，结果是买水果种子，还得自己种",
	}
	return nil, map[string]interface{}{"mistake": mistakes[rand.Intn(len(mistakes))]}, nil
}

// 53. 生成随机情侣间的搞笑拌嘴
type CoupleFightInput struct{}

func CoupleFight(ctx context.Context, req *mcp.CallToolRequest, input CoupleFightInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	fights := []string{
		"男：你为什么又买口红？女：你为什么又买游戏皮肤？男：口红你用得完吗？女：游戏皮肤你打得赢吗？",
		"女：你看那个女生穿的衣服好看吗？男：没注意。女：你为什么不看？男：看了怕你吃醋。女：你就是心里有鬼！",
		"男：今天吃什么？女：随便。男：吃火锅？女：太辣。男：吃中餐？女：没胃口。男：那你想吃什么？女：随便。",
		"女：我胖了吗？男：没有。女：你骗人，我都胖了五斤了。男：那有点胖。女：你居然说我胖！",
		"男：我们去看电影吧。女：好啊，看什么？男：看动作片。女：不要，我想看爱情片。男：那看喜剧片？女：都行。男：那就看动作片吧。女：你一点都不在乎我的感受！",
	}
	return nil, map[string]interface{}{"fight": fights[rand.Intn(len(fights))]}, nil
}

// 54. 生成随机朋友间的互怼日常
type FriendRoastInput struct{}

func FriendRoast(ctx context.Context, req *mcp.CallToolRequest, input FriendRoastInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	roasts := []string{
		"你这智商，怕是被门夹过吧，不然怎么会想出这种馊主意",
		"就你这颜值，拍照不P图都不好意思发朋友圈吧",
		"你这厨艺，做的饭狗都不吃，也就我给你面子才吃两口",
		"你这游戏打得，菜得抠脚，还好意思叫我带带你",
		"你这品味，穿的衣服像是从垃圾桶里捡来的，能不能换一套",
	}
	return nil, map[string]interface{}{"roast": roasts[rand.Intn(len(roasts))]}, nil
}

// 55. 生成随机老师的经典口头禅
type TeacherLineInput struct{}

func TeacherLine(ctx context.Context, req *mcp.CallToolRequest, input TeacherLineInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lines := []string{
		"这道题我讲最后一遍，听懂了吗？",
		"你们是我教过最差的一届学生",
		"看我干什么？看黑板！看黑板干什么？看书！",
		"体育老师今天有事，这节课上数学",
		"等你们上了大学就轻松了",
	}
	return nil, map[string]interface{}{"line": lines[rand.Intn(len(lines))]}, nil
}

// 56. 生成随机老板的经典口头禅
type BossLineInput struct{}

func BossLine(ctx context.Context, req *mcp.CallToolRequest, input BossLineInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lines := []string{
		"这个项目必须在明天早上之前完成，没完成的加班也要做完",
		"我不管你用什么方法，我只要结果",
		"年轻人，要多吃苦，多奋斗，不要怕累",
		"这个月的业绩怎么回事？再这样下去你们都得卷铺盖走人",
		"我当初创业的时候，比你们辛苦多了",
	}
	return nil, map[string]interface{}{"line": lines[rand.Intn(len(lines))]}, nil
}

// 57. 生成随机父母的经典唠叨
type ParentNagInput struct{}

func ParentNag(ctx context.Context, req *mcp.CallToolRequest, input ParentNagInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	nags := []string{
		"多穿点衣服，别感冒了，感冒了又得花钱看病",
		"别老玩手机，对眼睛不好，有空多看书",
		"早点睡觉，别熬夜，熬夜对身体不好",
		"多吃点饭，看你瘦的，风一吹就倒",
		"什么时候带个对象回来看看？你同学都结婚了",
	}
	return nil, map[string]interface{}{"nag": nags[rand.Intn(len(nags))]}, nil
}

// 58. 生成随机吃货的点菜纠结
type FoodOrderInput struct{}

func FoodOrder(ctx context.Context, req *mcp.CallToolRequest, input FoodOrderInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	orders := []string{
		"这个看起来好好吃，那个也不错，到底点哪个呢？",
		"我想吃辣的，但是又怕上火，怎么办？",
		"这个太贵了，那个好像不太好吃，好纠结",
		"要不点这个吧，不行，还是点那个吧，算了，还是点这个吧",
		"要不我们换一家吧，这家好像没有我想吃的",
	}
	return nil, map[string]interface{}{"order": orders[rand.Intn(len(orders))]}, nil
}

// 59. 生成随机打工人的周一综合征
type MondaySyndromeInput struct{}

func MondaySyndrome(ctx context.Context, req *mcp.CallToolRequest, input MondaySyndromeInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	syndromes := []string{
		"不想起床，不想上班，想继续睡觉，周一为什么要上班",
		"一想到周一要开会，就头疼，能不能请假不去",
		"周一的地铁怎么这么挤，挤得我怀疑人生",
		"周一的工作怎么这么多，感觉永远做不完",
		"一到周一就没精神，喝咖啡都没用，只想摸鱼",
	}
	return nil, map[string]interface{}{"syndrome": syndromes[rand.Intn(len(syndromes))]}, nil
}

// 60. 生成随机学生的考试前焦虑
type ExamAnxietyInput struct{}

func ExamAnxiety(ctx context.Context, req *mcp.CallToolRequest, input ExamAnxietyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	anxieties := []string{
		"还有好多知识点没复习，明天就要考试了，怎么办",
		"万一考砸了怎么办，爸妈会骂我的，老师也会失望的",
		"我现在一点都记不住，脑子一片空白，完了完了",
		"别人都复习得那么好，就我什么都不会，肯定考不好",
		"今晚肯定睡不着了，明天考试肯定没精神",
	}
	return nil, map[string]interface{}{"anxiety": anxieties[rand.Intn(len(anxieties))]}, nil
}

// 61. 生成随机网友的奇葩提问
type StrangeQuestionInput struct{}

func StrangeQuestion(ctx context.Context, req *mcp.CallToolRequest, input StrangeQuestionInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	questions := []string{
		"吃了安眠药再喝咖啡，会睡着还是睡不着？",
		"如果我把镜子打碎了，镜子里的我会疼吗？",
		"用充电宝给充电宝充电，能充满吗？",
		"秃头的人洗头，用洗发水还是洗面奶？",
		"如果我在做梦的时候说梦话，梦里的人能听到吗？",
	}
	return nil, map[string]interface{}{"question": questions[rand.Intn(len(questions))]}, nil
}

// 62. 生成随机猫咪的迷惑行为
type CatConfuseInput struct{}

func CatConfuse(ctx context.Context, req *mcp.CallToolRequest, input CatConfuseInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	behaviors := []string{
		"把猫砂盆里的猫砂扒出来，然后在地板上拉屎",
		"半夜在房间里跑酷，把东西都打翻",
		"主人在工作的时候，非要趴在键盘上，不让主人工作",
		"把桌子上的东西推下去，然后装作若无其事的样子",
		"害怕黄瓜，看到黄瓜就吓得跳起来",
	}
	return nil, map[string]interface{}{"behavior": behaviors[rand.Intn(len(behaviors))]}, nil
}

// 63. 生成随机狗狗的可爱行为
type DogCuteInput struct{}

func DogCute(ctx context.Context, req *mcp.CallToolRequest, input DogCuteInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	behaviors := []string{
		"主人回家时，摇着尾巴在门口等，还会叼来拖鞋",
		"听到主人说“散步”，就会兴奋地转圈，还会自己叼来牵引绳",
		"主人难过的时候，会安静地趴在主人身边，舔主人的手",
		"看到主人吃东西，就会用可怜的眼神看着主人，求投喂",
		"睡觉的时候会打呼噜，还会做梦蹬腿",
	}
	return nil, map[string]interface{}{"behavior": behaviors[rand.Intn(len(behaviors))]}, nil
}

// 64. 生成随机天气的奇葩现象
type StrangeWeatherInput struct{}

func StrangeWeather(ctx context.Context, req *mcp.CallToolRequest, input StrangeWeatherInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	weathers := []struct {
		Phenomenon string
		Place      string
	}{
		{"下鱼雨，天上掉下来好多小鱼", "澳大利亚"},
		{"下血雨，雨水呈现红色，像血一样", "印度"},
		{"下冰雹，冰雹有拳头那么大，砸坏了很多东西", "美国"},
		{"同时出现太阳和下雨，还出现了两道彩虹", "中国云南"},
		{"下青蛙雨，天上掉下来很多小青蛙", "英国"},
	}
	w := weathers[rand.Intn(len(weathers))]
	return nil, map[string]interface{}{
		"phenomenon": w.Phenomenon,
		"place":      w.Place,
	}, nil
}

// 65. 生成随机梦境的奇怪场景
type StrangeDreamInput struct{}

func StrangeDream(ctx context.Context, req *mcp.CallToolRequest, input StrangeDreamInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	dreams := []string{
		"自己变成了一只鸟，在天上飞，但是却不会降落",
		"在一个没有尽头的走廊里奔跑，后面有什么东西在追",
		"和明星一起吃饭，但是明星的脸一直在变",
		"自己在考试，但是题目都是看不懂的符号",
		"房子里的家具都活了过来，在和自己说话",
	}
	return nil, map[string]interface{}{"dream": dreams[rand.Intn(len(dreams))]}, nil
}

// 66. 生成随机童年的奇葩玩具
type ChildhoodToyInput struct{}

func ChildhoodToy(ctx context.Context, req *mcp.CallToolRequest, input ChildhoodToyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	toys := []struct {
		Name string
		Desc string
	}{
		{"弹珠", "圆圆的玻璃球，能在地上滚来滚去，还能和小朋友比赛"},
		{"跳房子格子", "用粉笔画在地上的格子，单脚双脚跳着玩"},
		{"铁皮青蛙", "上了发条就会跳的青蛙，铁皮做的，很吵但很有趣"},
		{"泡泡胶", "能吹成大泡泡的胶，有点臭，但能玩一下午"},
		{"东南西北", "用纸折的，能算命，还能和小朋友玩角色扮演"},
	}
	t := toys[rand.Intn(len(toys))]
	return nil, map[string]interface{}{
		"toy":  t.Name,
		"desc": t.Desc,
	}, nil
}

// 67. 生成随机小时候的奇葩零食
type ChildhoodSnackInput struct{}

func ChildhoodSnack(ctx context.Context, req *mcp.CallToolRequest, input ChildhoodSnackInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	snacks := []struct {
		Name string
		Desc string
	}{
		{"辣条", "红色的，很辣，一毛钱一根，偷偷在学校里吃"},
		{"大大泡泡糖", "能吹很大的泡泡，还有各种口味，包装上有卡通图案"},
		{"唐僧肉", "其实是萝卜干，但是名字很吸引人，吃起来甜甜的"},
		{"口哨糖", "能吹出声的糖，一边吃糖一边吹口哨，很得意"},
		{"冰棍", "用塑料袋装的，一毛钱一根，夏天吃很凉快"},
	}
	s := snacks[rand.Intn(len(snacks))]
	return nil, map[string]interface{}{
		"snack": s.Name,
		"desc":  s.Desc,
	}, nil
}

// 68. 生成随机打工人的午餐纠结
type LunchConfuseInput struct{}

func LunchConfuse(ctx context.Context, req *mcp.CallToolRequest, input LunchConfuseInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	confuses := []string{
		"今天中午吃什么？外卖还是出去吃？",
		"这家外卖昨天吃过了，那家好像不好吃，好纠结",
		"出去吃又要排队，外卖又要等很久，怎么办",
		"想吃点好的，但是又怕贵，还是省钱吧",
		"减肥期间，中午吃什么才不会胖呢",
	}
	return nil, map[string]interface{}{"confuse": confuses[rand.Intn(len(confuses))]}, nil
}

// 69. 生成随机网购时的好评文学
type GoodReviewInput struct{}

func GoodReview(ctx context.Context, req *mcp.CallToolRequest, input GoodReviewInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	reviews := []string{
		"这东西太好了，超出我的预期，必须五星好评，已经推荐给朋友了",
		"卖家服务态度很好，物流也很快，东西质量没话说，下次还会再来",
		"太喜欢了，和图片一模一样，没有色差，性价比很高",
		"包装很精致，还送了小礼物，太贴心了，必须好评",
		"用了一段时间才来评价，真的很好用，值得购买",
	}
	return nil, map[string]interface{}{"review": reviews[rand.Intn(len(reviews))]}, nil
}

// 70. 生成随机学生的逃课理由
type SkipClassReasonInput struct{}

func SkipClassReason(ctx context.Context, req *mcp.CallToolRequest, input SkipClassReasonInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	reasons := []string{
		"我生病了，头疼得厉害，去不了学校",
		"我家的猫生病了，我要带它去看医生",
		"我自行车坏了，去不了学校",
		"我闹钟没响，起来的时候已经上课了",
		"我亲戚来了，身体不舒服，想在家休息",
	}
	return nil, map[string]interface{}{"reason": reasons[rand.Intn(len(reasons))]}, nil
}

// 71. 生成随机打工人的离职理由
type ResignReasonInput struct{}

func ResignReason(ctx context.Context, req *mcp.CallToolRequest, input ResignReasonInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	reasons := []string{
		"我觉得这个工作不适合我，想换个环境",
		"我家里有事，需要回家处理，可能要很久",
		"我想继续深造，提升自己，所以要辞职",
		"我找到了一份更适合我的工作，薪资待遇也更好",
		"我身体不太好，想休息一段时间，调理一下身体",
	}
	return nil, map[string]interface{}{"reason": reasons[rand.Intn(len(reasons))]}, nil
}

// 72. 生成随机网友的神评论
type GodCommentInput struct{}

func GodComment(ctx context.Context, req *mcp.CallToolRequest, input GodCommentInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	comments := []string{
		"看了你的视频，我终于知道我为什么单身了，因为我没你这么优秀",
		"这个操作，秀得我头皮发麻",
		"我以为是个青铜，没想到是个王者",
		"建议直接出道，我第一个投票",
		"别人笑我太疯癫，我笑他人看不穿",
	}
	return nil, map[string]interface{}{"comment": comments[rand.Intn(len(comments))]}, nil
}

// 73. 生成随机开车时的搞笑经历
type DrivingStoryInput struct{}

func DrivingStory(ctx context.Context, req *mcp.CallToolRequest, input DrivingStoryInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	stories := []string{
		"开车的时候，导航说“前方有急转弯，请减速”，结果我减速了，后面的车以为我要停车，差点追尾",
		"在停车场找不到自己的车了，绕了半个小时才找到，原来就在入口旁边",
		"开车的时候，一只鸟飞到了挡风玻璃上，吓得我一脚刹车，后面的车喇叭响个不停",
		"加油的时候，忘了熄火，加油站的工作人员吓得赶紧让我熄火",
		"倒车的时候，没看到后面的电线杆，砰的一声撞上去了，车屁股凹了一块",
	}
	return nil, map[string]interface{}{"story": stories[rand.Intn(len(stories))]}, nil
}

// 74. 生成随机做饭时的翻车现场
type CookingFailInput struct{}

func CookingFail(ctx context.Context, req *mcp.CallToolRequest, input CookingFailInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	fails := []string{
		"想煎个荷包蛋，结果鸡蛋煎糊了，锅也黑了，还差点把厨房点着",
		"煮面条的时候，忘了关火，水都煮干了，面条变成了炭",
		"想做个蛋糕，结果蛋糕没发起来，变成了饼，还特别硬",
		"炒青菜的时候，盐放多了，咸得没法吃，只能倒掉",
		"炖排骨的时候，忘了盖锅盖，汤都炖没了，排骨也炖老了",
	}
	return nil, map[string]interface{}{"fail": fails[rand.Intn(len(fails))]}, nil
}

// 75. 生成随机自拍时的搞笑姿势
type SelfiePoseInput struct{}

func SelfiePose(ctx context.Context, req *mcp.CallToolRequest, input SelfiePoseInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	poses := []string{
		"用手比个心，放在脸旁边，眼睛瞪得大大的",
		"把头发撩起来，露出额头，假装很酷的样子",
		"用手捂住一只眼睛，另一只眼睛眨一下",
		"嘴巴嘟起来，像个小金鱼，再配上无辜的眼神",
		"跳起来拍，结果拍成了表情包",
	}
	return nil, map[string]interface{}{"pose": poses[rand.Intn(len(poses))]}, nil
}

// 76. 生成随机聚会时的游戏
type PartyGameInput struct{}

func PartyGame(ctx context.Context, req *mcp.CallToolRequest, input PartyGameInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	games := []struct {
		Name string
		Rule string
	}{
		{"真心话大冒险", "通过石头剪刀布决定输赢，输的人选择真心话或大冒险"},
		{"狼人杀", "多人游戏，有狼人、村民、预言家等角色，通过发言找出狼人"},
		{"谁是卧底", "每个人拿到一个词语，其中一个人拿到的是卧底词，通过描述找出卧底"},
		{"国王游戏", "抽牌决定国王，国王可以命令其他人做事情"},
		{"你画我猜", "一个人画，其他人猜画的是什么，看谁猜得快"},
	}
	g := games[rand.Intn(len(games))]
	return nil, map[string]interface{}{
		"game": g.Name,
		"rule": g.Rule,
	}, nil
}

// 77. 生成随机KTV必点歌曲
type KTVSongInput struct{}

func KTVSong(ctx context.Context, req *mcp.CallToolRequest, input KTVSongInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	songs := []struct {
		Name   string
		Singer string
	}{
		{"《死了都要爱》", "信乐团"},
		{"《王妃》", "萧敬腾"},
		{"《小苹果》", "筷子兄弟"},
		{"《江南》", "林俊杰"},
		{"《后来》", "刘若英"},
	}
	s := songs[rand.Intn(len(songs))]
	return nil, map[string]interface{}{
		"song":   s.Name,
		"singer": s.Singer,
	}, nil
}

// 78. 生成随机健身时的摸鱼行为
type FitnessSlackInput struct{}

func FitnessSlack(ctx context.Context, req *mcp.CallToolRequest, input FitnessSlackInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	behaviors := []string{
		"在跑步机上慢慢走，假装在跑步，其实在看手机",
		"举哑铃的时候，只举几下就放下，然后去喝水休息",
		"别人在做高强度训练，自己在旁边做拉伸，拉伸了半个小时",
		"去健身房换了衣服，然后在休息区坐着玩手机，玩了一个小时就走了",
		"假装去洗手间，其实在里面刷短视频，刷了十几分钟",
	}
	return nil, map[string]interface{}{"behavior": behaviors[rand.Intn(len(behaviors))]}, nil
}

// 79. 生成随机网购时的省钱技巧
type ShoppingSaveInput struct{}

func ShoppingSave(ctx context.Context, req *mcp.CallToolRequest, input ShoppingSaveInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	skills := []string{
		"在商品详情页领优惠券，很多优惠券不领就看不到",
		"在节日的时候买东西，比如双十一、618，折扣比较大",
		"货比三家，多看看不同的店铺，找性价比最高的",
		"关注店铺的直播，直播的时候经常有秒杀活动",
		"把想买的东西加入购物车，等降价了再买",
	}
	return nil, map[string]interface{}{"skill": skills[rand.Intn(len(skills))]}, nil
}

// 80. 生成随机职场中的潜规则
type WorkplaceRuleInput struct{}

func WorkplaceRule(ctx context.Context, req *mcp.CallToolRequest, input WorkplaceRuleInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	rules := []string{
		"领导说“随便看看”，其实是想让你认真看，并且给出意见",
		"同事说“下次请你吃饭”，大部分时候都是客套话，别当真",
		"开会的时候，领导最后发言，不要抢在领导前面说太多",
		"不要在背后议论同事和领导，坏话总会传到他们耳朵里",
		"收到消息要及时回复，哪怕只是回复一个“好的”",
	}
	return nil, map[string]interface{}{"rule": rules[rand.Intn(len(rules))]}, nil
}

// 81. 生成随机校园里的奇葩规定
type SchoolRuleInput struct{}

func SchoolRule(ctx context.Context, req *mcp.CallToolRequest, input SchoolRuleInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	rules := []string{
		"女生不能留长发，必须剪短发，说是为了方便打理",
		"走路必须走直线，不能并排走，说是为了保持秩序",
		"不许带零食进校园，发现了就没收，还要通报批评",
		"周一到周五必须穿校服，哪怕是夏天也不能穿自己的衣服",
		"晚上10点必须关灯睡觉，不许玩手机，老师会查房",
	}
	return nil, map[string]interface{}{"rule": rules[rand.Intn(len(rules))]}, nil
}

// 82. 生成随机恋爱中的甜蜜小事
type LoveSweetInput struct{}

func LoveSweet(ctx context.Context, req *mcp.CallToolRequest, input LoveSweetInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	sweets := []string{
		"下雨的时候，他把伞都倾向我这边，自己半边身子都淋湿了",
		"我随口说想吃某家的蛋糕，他跑了很远的路给我买回来",
		"晚上睡觉的时候，他会把我抱得很紧，怕我踢被子",
		"我来例假的时候，他会给我煮红糖姜茶，还会给我揉肚子",
		"他记得我所有的喜好，知道我不吃香菜，喜欢喝奶茶三分甜",
	}
	return nil, map[string]interface{}{"sweet": sweets[rand.Intn(len(sweets))]}, nil
}

// 83. 生成随机朋友间的暖心瞬间
type FriendWarmInput struct{}

func FriendWarm(ctx context.Context, req *mcp.CallToolRequest, input FriendWarmInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	warms := []string{
		"我难过的时候，她二话不说就过来陪我，听我吐槽了一晚上",
		"我没钱的时候，他主动借给我，还说不急着还",
		"我生病的时候，她给我送来了药和粥，还帮我打扫了房间",
		"我失恋的时候，他拉着我去吃火锅，陪我喝酒，说天涯何处无芳草",
		"我找工作不顺利的时候，他一直在鼓励我，还帮我修改简历",
	}
	return nil, map[string]interface{}{"warm": warms[rand.Intn(len(warms))]}, nil
}

// 84. 生成随机家人间的温馨时刻
type FamilyWarmInput struct{}

func FamilyWarm(ctx context.Context, req *mcp.CallToolRequest, input FamilyWarmInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	warms := []string{
		"我回家的时候，妈妈早就做好了一桌子我爱吃的菜",
		"我生病的时候，爸爸半夜起来给我倒水，还一直守在我床边",
		"我工作不顺心的时候，爷爷安慰我说，没关系，慢慢来，家里永远是你的后盾",
		"我过生日的时候，全家人都记得，还一起给我唱生日歌",
		"我出门的时候，奶奶一直叮嘱我要注意安全，早点回家",
	}
	return nil, map[string]interface{}{"warm": warms[rand.Intn(len(warms))]}, nil
}

// 85. 生成随机旅行中的暖心经历
type TravelWarmInput struct{}

func TravelWarm(ctx context.Context, req *mcp.CallToolRequest, input TravelWarmInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	warms := []string{
		"在陌生的城市迷路了，一个老奶奶主动带我找到了目的地",
		"在火车上，旁边的大叔看到我没带吃的，给了我一个面包",
		"在景区排队的时候，前面的小姐姐看到我很累，让我站到她前面",
		"在酒店住的时候，服务员看到我感冒了，主动给我送来了感冒药",
		"在爬山的时候，看到一个小朋友摔倒了，一个陌生人赶紧跑过去把他扶起来",
	}
	return nil, map[string]interface{}{"warm": warms[rand.Intn(len(warms))]}, nil
}

// 86. 生成随机生活中的小确幸
type LittleHappinessInput struct{}

func LittleHappiness(ctx context.Context, req *mcp.CallToolRequest, input LittleHappinessInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	happiness := []string{
		"早上醒来，发现今天是周末，可以不用上班",
		"走在路上，捡到了一块钱，虽然不多，但很开心",
		"买咖啡的时候，店员多给了我一个小饼干",
		"下雨的时候，刚好带了伞，而别人在淋雨",
		"晚上睡觉的时候，发现被窝里很暖和",
	}
	return nil, map[string]interface{}{"happiness": happiness[rand.Intn(len(happiness))]}, nil
}

// 87. 生成随机动物的可爱瞬间
type AnimalCuteInput struct{}

func AnimalCute(ctx context.Context, req *mcp.CallToolRequest, input AnimalCuteInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	cutes := []string{
		"小猫咪蜷缩在阳光下睡觉，还打了个小呼噜",
		"小狗看到主人回家，兴奋地摇着尾巴，还在地上打滚",
		"小熊猫抱着竹子啃，吃得满脸都是",
		"小兔子三瓣嘴一动一动地吃胡萝卜，耳朵还时不时动一下",
		"小松鼠在树枝上跳来跳去，还抱着一颗松果",
	}
	return nil, map[string]interface{}{"cute": cutes[rand.Intn(len(cutes))]}, nil
}

// 88. 生成随机自然的美丽景色
type NatureBeautyInput struct{}

func NatureBeauty(ctx context.Context, req *mcp.CallToolRequest, input NatureBeautyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	beauties := []struct {
		Scene string
		Time  string
	}{
		{"海边的日落，太阳把天空染成了红色和橙色", "傍晚"},
		{"山顶的云海，云雾缭绕，像仙境一样", "清晨"},
		{"森林里的小溪，溪水清澈见底，还有小鱼在游", "中午"},
		{"田野里的油菜花，一片金黄，还引来很多蜜蜂", "春天"},
		{"冬天的雪景，大地一片洁白，树枝上挂满了雪花", "冬天"},
	}
	b := beauties[rand.Intn(len(beauties))]
	return nil, map[string]interface{}{
		"scene": b.Scene,
		"time":  b.Time,
	}, nil
}

// 89. 生成随机城市的夜景
type CityNightInput struct{}

func CityNight(ctx context.Context, req *mcp.CallToolRequest, input CityNightInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	nights := []struct {
		City  string
		Scene string
	}{
		{"上海", "外滩的夜景，灯光璀璨，东方明珠塔格外显眼"},
		{"北京", "天安门广场的夜景，庄严肃穆，灯火通明"},
		{"广州", "珠江的夜景，游船穿梭，两岸的灯光倒映在水里"},
		{"成都", "锦里的夜景，古色古香的建筑配上红灯笼，很有韵味"},
		{"西安", "大唐不夜城的夜景，仿佛穿越回了唐朝"},
	}
	n := nights[rand.Intn(len(nights))]
	return nil, map[string]interface{}{
		"city":  n.City,
		"scene": n.Scene,
	}, nil
}

// 90. 生成随机季节的美好事物
type SeasonBeautyInput struct{}

func SeasonBeauty(ctx context.Context, req *mcp.CallToolRequest, input SeasonBeautyInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	beauties := []struct {
		Season string
		Thing  string
	}{
		{"春天", "公园里的樱花，粉粉嫩嫩的，随风飘落像雪花"},
		{"夏天", "夜晚的萤火虫，一闪一闪的，像天上的星星"},
		{"秋天", "路边的枫叶，红红的，像燃烧的火焰"},
		{"冬天", "窗外的雪花，飘飘扬扬的，把世界变得洁白"},
	}
	b := beauties[rand.Intn(len(beauties))]
	return nil, map[string]interface{}{
		"season": b.Season,
		"thing":  b.Thing,
	}, nil
}

// 91. 生成随机读书时的感悟
type ReadingFeelingInput struct{}

func ReadingFeeling(ctx context.Context, req *mcp.CallToolRequest, input ReadingFeelingInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	feelings := []string{
		"读一本书，就像认识一个新朋友，能学到很多东西",
		"书里的世界很精彩，能让人忘记现实中的烦恼",
		"有些书，第一次读和第二次读，会有不同的感受",
		"读书能让人变得安静，也能让人变得强大",
		"每本书都有它的灵魂，需要用心去感受",
	}
	return nil, map[string]interface{}{"feeling": feelings[rand.Intn(len(feelings))]}, nil
}

// 92. 生成随机电影中的经典台词
type MovieLineInput struct{}

func MovieLine(ctx context.Context, req *mcp.CallToolRequest, input MovieLineInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lines := []struct {
		Line  string
		Movie string
	}{
		{"曾经有一份真挚的爱情摆在我的面前，但是我没有珍惜...", "《大话西游》"},
		{"人生就像一盒巧克力，你永远不知道下一块会是什么味道。", "《阿甘正传》"},
		{"我猜中了开头，却猜不中这结局。", "《大话西游》"},
		{"如果不能骄傲地活着，我选择死亡。", "《霸王别姬》"},
		{"世界上有一种鸟是没有脚的，它只能够一直飞...", "《阿飞正传》"},
	}
	l := lines[rand.Intn(len(lines))]
	return nil, map[string]interface{}{
		"line":  l.Line,
		"movie": l.Movie,
	}, nil
}

// 93. 生成随机歌曲中的经典歌词
type SongLyricInput struct{}

func SongLyric(ctx context.Context, req *mcp.CallToolRequest, input SongLyricInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lyrics := []struct {
		Lyric string
		Song  string
	}{
		{"后来，终于在眼泪中明白，有些人，一旦错过就不在。", "《后来》"},
		{"听妈妈的话，别让她受伤，想快快长大，才能保护她。", "《听妈妈的话》"},
		{"我可以抱你吗爱人，让我在你肩膀哭泣。", "《我可以抱你吗》"},
		{"十年之前，我不认识你，你不属于我，我们还是一样...", "《十年》"},
		{"阳光总在风雨后，请相信有彩虹。", "《阳光总在风雨后》"},
	}
	l := lyrics[rand.Intn(len(lyrics))]
	return nil, map[string]interface{}{
		"lyric": l.Lyric,
		"song":  l.Song,
	}, nil
}

// 94. 生成随机游戏中的经典台词
type GameLineInput struct{}

func GameLine(ctx context.Context, req *mcp.CallToolRequest, input GameLineInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lines := []struct {
		Line string
		Game string
	}{
		{"为了部落！", "《魔兽世界》"},
		{"德玛西亚！", "《英雄联盟》"},
		{"安息吧，我的朋友。", "《暗黑破坏神》"},
		{"我是要成为海贼王的男人！", "《海贼王：无尽世界》"},
		{"你已经死了。", "《北斗神拳》游戏版"},
	}
	l := lines[rand.Intn(len(lines))]
	return nil, map[string]interface{}{
		"line": l.Line,
		"game": l.Game,
	}, nil
}

// 95. 生成随机动漫中的经典台词
type AnimeLineInput struct{}

func AnimeLine(ctx context.Context, req *mcp.CallToolRequest, input AnimeLineInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	lines := []struct {
		Line  string
		Anime string
	}{
		{"我要代表月亮消灭你！", "《美少女战士》"},
		{"真相只有一个！", "《名侦探柯南》"},
		{"我是要成为火影的男人！", "《火影忍者》"},
		{"海贼王，我当定了！", "《海贼王》"},
		{"你还差得远呢！", "《网球王子》"},
	}
	l := lines[rand.Intn(len(lines))]
	return nil, map[string]interface{}{
		"line":  l.Line,
		"anime": l.Anime,
	}, nil
}

// 96. 生成随机历史人物的经典名言
type HistoricalQuoteInput struct{}

func HistoricalQuote(ctx context.Context, req *mcp.CallToolRequest, input HistoricalQuoteInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	quotes := []struct {
		Quote  string
		Person string
	}{
		{"人生自古谁无死，留取丹心照汗青。", "文天祥"},
		{"三人行，必有我师焉。", "孔子"},
		{"天生我材必有用，千金散尽还复来。", "李白"},
		{"苟利国家生死以，岂因祸福避趋之。", "林则徐"},
		{"先天下之忧而忧，后天下之乐而乐。", "范仲淹"},
	}
	q := quotes[rand.Intn(len(quotes))]
	return nil, map[string]interface{}{
		"quote":  q.Quote,
		"person": q.Person,
	}, nil
}

// 97. 生成随机名人的励志名言
type CelebrityQuoteInput struct{}

func CelebrityQuote(ctx context.Context, req *mcp.CallToolRequest, input CelebrityQuoteInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	quotes := []struct {
		Quote     string
		Celebrity string
	}{
		{"成功不是偶然的，而是必然的。", "马云"},
		{"天才是百分之一的灵感加上百分之九十九的汗水。", "爱迪生"},
		{"不想当将军的士兵不是好士兵。", "拿破仑"},
		{"生命就像一盒巧克力，结果往往出人意料。", " Forrest Gump"},
		{"我们最大的光荣不在于永不失败，而在于每次跌倒后都能爬起来。", "丘吉尔"},
	}
	q := quotes[rand.Intn(len(quotes))]
	return nil, map[string]interface{}{
		"quote":     q.Quote,
		"celebrity": q.Celebrity,
	}, nil
}

// 98. 生成随机生活中的励志瞬间
type InspirationalMomentInput struct{}

func InspirationalMoment(ctx context.Context, req *mcp.CallToolRequest, input InspirationalMomentInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	moments := []string{
		"坚持跑步一个月，终于瘦了五斤，感觉自己很棒",
		"努力学习了很久，终于通过了考试，付出有了回报",
		"第一次做饭，虽然有点难吃，但还是很开心，因为是自己做的",
		"克服了自己的恐惧，第一次蹦极，感觉很刺激，也很有成就感",
		"一直想做的事情，终于鼓起勇气去做了，不管结果如何，都不后悔",
	}
	return nil, map[string]interface{}{"moment": moments[rand.Intn(len(moments))]}, nil
}

// 99. 生成随机未来的小目标
type FutureGoalInput struct{}

func FutureGoal(ctx context.Context, req *mcp.CallToolRequest, input FutureGoalInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	goals := []string{
		"下个月瘦十斤，每天坚持运动一小时",
		"今年读完20本书，每个月至少读一本",
		"学会做一道拿手菜，给家人露一手",
		"存够钱，去一次西藏，看看那里的蓝天白云",
		"学一门新技能，比如画画或者弹吉他",
	}
	return nil, map[string]interface{}{"goal": goals[rand.Intn(len(goals))]}, nil
}

// 100. 生成随机人生的小感悟
type LifeFeelingInput struct{}

func LifeFeeling(ctx context.Context, req *mcp.CallToolRequest, input LifeFeelingInput) (*mcp.CallToolResult, map[string]interface{}, error) {
	feelings := []string{
		"人生就像一场旅行，重要的不是目的地，而是沿途的风景",
		"幸福其实很简单，就是身边有爱的人，有喜欢的事",
		"不要太在意别人的眼光，做好自己就好",
		"珍惜当下，因为明天和意外不知道哪个会先来",
		"努力不一定会成功，但不努力一定不会成功",
	}
	return nil, map[string]interface{}{"feeling": feelings[rand.Intn(len(feelings))]}, nil
}

func main() {
	// 初始化随机数生成器（确保随机结果不同）
	rand.Seed(time.Now().UnixNano())

	// 创建MCP服务器实例
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "fun-tools-collection",
			Version: "v1.0.0",
		},
		nil,
	)

	// 注册所有工具函数（共100个）
	// 1. 彩虹屁生成器
	mcp.AddTool[RainbowFartInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "rainbowFart", Description: "生成一句有趣的彩虹屁夸赞"},
		RainbowFart,
	)

	// 2. 猜拳游戏
	mcp.AddTool[RockPaperScissorsInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "rockPaperScissors", Description: "与AI进行猜拳游戏（rock/paper/scissors）"},
		RockPaperScissors,
	)

	// 3. 中二台词生成
	mcp.AddTool[ChuuniLineInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "chuuniLine", Description: "生成一句随机中二台词"},
		ChuuniLine,
	)

	// 4. 冷门电影推荐
	mcp.AddTool[ObscureMovieInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "obscureMovie", Description: "推荐一部冷门但优质的电影"},
		ObscureMovie,
	)

	// 5. 城市小众景点推荐
	mcp.AddTool[HiddenSpotInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "hiddenSpot", Description: "推荐指定城市的小众景点"},
		HiddenSpot,
	)

	// 6. 无用但有趣的知识
	mcp.AddTool[UselessFactInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "uselessFact", Description: "提供一个无用但有趣的冷知识"},
		UselessFact,
	)

	// 7. 早餐搭配推荐
	mcp.AddTool[BreakfastComboInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "breakfastCombo", Description: "生成一份随机早餐搭配方案"},
		BreakfastCombo,
	)

	// 8. Emoji故事生成
	mcp.AddTool[EmojiStoryInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "emojiStory", Description: "用3个emoji组成一个小故事"},
		EmojiStory,
	)

	// 9. 宠物中二名字
	mcp.AddTool[PetChuuniNameInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "petChuuniName", Description: "给宠物起一个中二风格的名字"},
		PetChuuniName,
	)

	// 11. 朋友圈文案生成
	mcp.AddTool[MomentsCaptionInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "momentsCaption", Description: "根据心情生成朋友圈文案"},
		MomentsCaption,
	)

	// 12. 解压小方法
	mcp.AddTool[StressReliefInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "stressRelief", Description: "提供一个有趣的解压小方法"},
		StressRelief,
	)

	// 13. 睡前小故事
	mcp.AddTool[BedtimeStoryInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "bedtimeStory", Description: "生成一句简短的睡前小故事"},
		BedtimeStory,
	)

	// 14. 冷门爱好推荐
	mcp.AddTool[ObscureHobbyInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "obscureHobby", Description: "推荐一种冷门但有趣的爱好"},
		ObscureHobby,
	)

	// 15. 咖啡拉花幻想
	mcp.AddTool[CoffeeArtInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "coffeeArt", Description: "生成一个幻想中的咖啡拉花图案"},
		CoffeeArt,
	)

	// 16. 方言打招呼
	mcp.AddTool[DialectGreetingInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "dialectGreeting", Description: "用随机方言打招呼（附带翻译）"},
		DialectGreeting,
	)

	// 17. 网络热梗变体
	mcp.AddTool[MemeVariantInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "memeVariant", Description: "生成网络热梗的搞笑变体"},
		MemeVariant,
	)

	// 18. 奇葩零食搭配
	mcp.AddTool[WeirdSnackComboInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "weirdSnackCombo", Description: "推荐一种奇葩的零食搭配"},
		WeirdSnackCombo,
	)

	// 19. 做梦素材
	mcp.AddTool[DreamMaterialInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "dreamMaterial", Description: "提供一个有趣的做梦素材"},
		DreamMaterial,
	)

	// 20. 摸鱼黑话
	mcp.AddTool[FishLanguageInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "fishLanguage", Description: "生成老板听不懂的摸鱼黑话"},
		FishLanguage,
	)

	// 21. 天气梗生成
	mcp.AddTool[WeatherMemeInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "weatherMeme", Description: "根据天气生成搞笑梗"},
		WeatherMeme,
	)

	// 22. 洗澡歌曲推荐
	mcp.AddTool[ShowerSongInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "showerSong", Description: "推荐一首洗澡时适合唱的歌"},
		ShowerSong,
	)

	// 23. 植物吐槽
	mcp.AddTool[PlantRoastInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "plantRoast", Description: "生成对指定植物的搞笑吐槽"},
		PlantRoast,
	)

	// 24. 奇怪的节日
	mcp.AddTool[WeirdHolidayInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "weirdHoliday", Description: "生成一个奇怪的节日及规则"},
		WeirdHoliday,
	)

	// 25. 宠物内心戏
	mcp.AddTool[PetThoughtInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "petThought", Description: "模拟指定宠物的内心想法"},
		PetThought,
	)

	// 26. 古文版流行语
	mcp.AddTool[ClassicMemeInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "classicMeme", Description: "将网络流行语翻译成古文"},
		ClassicMeme,
	)

	// 27. 奇怪的解压玩具
	mcp.AddTool[WeirdFidgetToyInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "weirdFidgetToy", Description: "推荐一种奇怪的解压玩具"},
		WeirdFidgetToy,
	)

	// 28. 失眠胡思乱想
	mcp.AddTool[InsomniaThoughtInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "insomniaThought", Description: "生成失眠时的搞笑胡思乱想"},
		InsomniaThought,
	)

	// 29. 情侣幼稚游戏
	mcp.AddTool[CuteCoupleGameInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "cuteCoupleGame", Description: "推荐情侣间的幼稚小游戏"},
		CuteCoupleGame,
	)

	// 30. 外卖备注骚话
	mcp.AddTool[TakeawayNoteInput, map[string]interface{}](
		server,
		&mcp.Tool{Name: "takeawayNote", Description: "生成外卖备注的搞笑骚话"},
		TakeawayNote,
	)

	// 31-100. 剩余工具注册（格式与上面一致，此处省略重复代码）
	// 31. 职场摸鱼借口
	mcp.AddTool[WorkSlackExcuseInput, map[string]interface{}](server, &mcp.Tool{Name: "workSlackExcuse", Description: "生成职场摸鱼的借口"}, WorkSlackExcuse)
	// 32. 网友抬杠语录
	mcp.AddTool[NetizenArgueInput, map[string]interface{}](server, &mcp.Tool{Name: "netizenArgue", Description: "生成网友抬杠语录"}, NetizenArgue)
	// 33. 减肥失败理由
	mcp.AddTool[DietFailReasonInput, map[string]interface{}](server, &mcp.Tool{Name: "dietFailReason", Description: "生成减肥失败的搞笑理由"}, DietFailReason)
	// 34. 朋友圈分组名称
	mcp.AddTool[MomentsGroupInput, map[string]interface{}](server, &mcp.Tool{Name: "momentsGroup", Description: "生成朋友圈分组名称"}, MomentsGroup)
	// 35. 网购收货名
	mcp.AddTool[ShoppingNameInput, map[string]interface{}](server, &mcp.Tool{Name: "shoppingName", Description: "生成搞笑的网购收货名"}, ShoppingName)
	// 36. 堵车内心OS
	mcp.AddTool[TrafficJamOSInput, map[string]interface{}](server, &mcp.Tool{Name: "trafficJamOS", Description: "生成堵车时的内心OS"}, TrafficJamOS)
	// 37. 考试前迷信行为
	mcp.AddTool[ExamSuperstitionInput, map[string]interface{}](server, &mcp.Tool{Name: "examSuperstition", Description: "生成考试前的迷信行为"}, ExamSuperstition)
	// 38. 游戏嘴强语录
	mcp.AddTool[GameTrashTalkInput, map[string]interface{}](server, &mcp.Tool{Name: "gameTrashTalk", Description: "生成打游戏时的嘴强语录"}, GameTrashTalk)
	// 39. 失眠自我安慰
	mcp.AddTool[InsomniaComfortInput, map[string]interface{}](server, &mcp.Tool{Name: "insomniaComfort", Description: "生成失眠时的自我安慰"}, InsomniaComfort)
	// 40. 被催婚反击
	mcp.AddTool[MarriageUrgeReplyInput, map[string]interface{}](server, &mcp.Tool{Name: "marriageUrgeReply", Description: "生成被催婚时的反击话术"}, MarriageUrgeReply)
	// 41. 老板画饼语录
	mcp.AddTool[BossPromiseInput, map[string]interface{}](server, &mcp.Tool{Name: "bossPromise", Description: "生成老板画饼的经典语录"}, BossPromise)
	// 42. 网购差评文学
	mcp.AddTool[BadReviewInput, map[string]interface{}](server, &mcp.Tool{Name: "badReview", Description: "生成搞笑的网购差评"}, BadReview)
	// 43. 减肥自我欺骗
	mcp.AddTool[DietCheatInput, map[string]interface{}](server, &mcp.Tool{Name: "dietCheat", Description: "生成减肥时的自我欺骗话术"}, DietCheat)
	// 44. 学生时代借口
	mcp.AddTool[StudentExcuseInput, map[string]interface{}](server, &mcp.Tool{Name: "studentExcuse", Description: "生成学生时代的经典借口"}, StudentExcuse)
	// 45. 家长群戏精发言
	mcp.AddTool[ParentGroupInput, map[string]interface{}](server, &mcp.Tool{Name: "parentGroup", Description: "生成家长群里的戏精发言"}, ParentGroup)
	// 46. 打工人周末计划
	mcp.AddTool[WeekendPlanInput, map[string]interface{}](server, &mcp.Tool{Name: "weekendPlan", Description: "生成打工人的周末计划"}, WeekendPlan)
	// 47. 吃货人生感悟
	mcp.AddTool[FoodieFeelingInput, map[string]interface{}](server, &mcp.Tool{Name: "foodieFeeling", Description: "生成吃货的人生感悟"}, FoodieFeeling)
	// 48. 深夜emo文案
	mcp.AddTool[LateNightEmoInput, map[string]interface{}](server, &mcp.Tool{Name: "lateNightEmo", Description: "生成深夜emo的朋友圈文案"}, LateNightEmo)
	// 49. 网友迷惑行为
	mcp.AddTool[ConfusedBehaviorInput, map[string]interface{}](server, &mcp.Tool{Name: "confusedBehavior", Description: "生成网友的迷惑行为"}, ConfusedBehavior)
	// 50. 摸鱼小技巧
	mcp.AddTool[SlackSkillInput, map[string]interface{}](server, &mcp.Tool{Name: "slackSkill", Description: "生成打工人的摸鱼小技巧"}, SlackSkill)
	// 51. 旅行奇葩经历
	mcp.AddTool[TravelStoryInput, map[string]interface{}](server, &mcp.Tool{Name: "travelStory", Description: "生成旅行中的奇葩经历"}, TravelStory)
	// 52. 网购搞笑误会
	mcp.AddTool[ShoppingMistakeInput, map[string]interface{}](server, &mcp.Tool{Name: "shoppingMistake", Description: "生成网购时的搞笑误会"}, ShoppingMistake)
	// 53. 情侣搞笑拌嘴
	mcp.AddTool[CoupleFightInput, map[string]interface{}](server, &mcp.Tool{Name: "coupleFight", Description: "生成情侣间的搞笑拌嘴"}, CoupleFight)
	// 54. 朋友互怼日常
	mcp.AddTool[FriendRoastInput, map[string]interface{}](server, &mcp.Tool{Name: "friendRoast", Description: "生成朋友间的互怼日常"}, FriendRoast)
	// 55. 老师口头禅
	mcp.AddTool[TeacherLineInput, map[string]interface{}](server, &mcp.Tool{Name: "teacherLine", Description: "生成老师的经典口头禅"}, TeacherLine)
	// 56. 老板口头禅
	mcp.AddTool[BossLineInput, map[string]interface{}](server, &mcp.Tool{Name: "bossLine", Description: "生成老板的经典口头禅"}, BossLine)
	// 57. 父母唠叨
	mcp.AddTool[ParentNagInput, map[string]interface{}](server, &mcp.Tool{Name: "parentNag", Description: "生成父母的经典唠叨"}, ParentNag)
	// 58. 吃货点菜纠结
	mcp.AddTool[FoodOrderInput, map[string]interface{}](server, &mcp.Tool{Name: "foodOrder", Description: "生成吃货的点菜纠结"}, FoodOrder)
	// 59. 周一综合征
	mcp.AddTool[MondaySyndromeInput, map[string]interface{}](server, &mcp.Tool{Name: "mondaySyndrome", Description: "生成打工人的周一综合征"}, MondaySyndrome)
	// 60. 考试前焦虑
	mcp.AddTool[ExamAnxietyInput, map[string]interface{}](server, &mcp.Tool{Name: "examAnxiety", Description: "生成学生的考试前焦虑"}, ExamAnxiety)
	// 61. 网友奇葩提问
	mcp.AddTool[StrangeQuestionInput, map[string]interface{}](server, &mcp.Tool{Name: "strangeQuestion", Description: "生成网友的奇葩提问"}, StrangeQuestion)
	// 62. 猫咪迷惑行为
	mcp.AddTool[CatConfuseInput, map[string]interface{}](server, &mcp.Tool{Name: "catConfuse", Description: "生成猫咪的迷惑行为"}, CatConfuse)
	// 63. 狗狗可爱行为
	mcp.AddTool[DogCuteInput, map[string]interface{}](server, &mcp.Tool{Name: "dogCute", Description: "生成狗狗的可爱行为"}, DogCute)
	// 64. 天气奇葩现象
	mcp.AddTool[StrangeWeatherInput, map[string]interface{}](server, &mcp.Tool{Name: "strangeWeather", Description: "生成天气的奇葩现象"}, StrangeWeather)
	// 65. 梦境奇怪场景
	mcp.AddTool[StrangeDreamInput, map[string]interface{}](server, &mcp.Tool{Name: "strangeDream", Description: "生成梦境的奇怪场景"}, StrangeDream)
	// 66. 童年奇葩玩具
	mcp.AddTool[ChildhoodToyInput, map[string]interface{}](server, &mcp.Tool{Name: "childhoodToy", Description: "生成童年的奇葩玩具"}, ChildhoodToy)
	// 67. 童年奇葩零食
	mcp.AddTool[ChildhoodSnackInput, map[string]interface{}](server, &mcp.Tool{Name: "childhoodSnack", Description: "生成童年的奇葩零食"}, ChildhoodSnack)
	// 68. 午餐纠结
	mcp.AddTool[LunchConfuseInput, map[string]interface{}](server, &mcp.Tool{Name: "lunchConfuse", Description: "生成打工人的午餐纠结"}, LunchConfuse)
	// 69. 网购好评文学
	mcp.AddTool[GoodReviewInput, map[string]interface{}](server, &mcp.Tool{Name: "goodReview", Description: "生成搞笑的网购好评"}, GoodReview)
	// 70. 学生逃课理由
	mcp.AddTool[SkipClassReasonInput, map[string]interface{}](server, &mcp.Tool{Name: "skipClassReason", Description: "生成学生的逃课理由"}, SkipClassReason)
	// 71. 打工人离职理由
	mcp.AddTool[ResignReasonInput, map[string]interface{}](server, &mcp.Tool{Name: "resignReason", Description: "生成打工人的离职理由"}, ResignReason)
	// 72. 网友神评论
	mcp.AddTool[GodCommentInput, map[string]interface{}](server, &mcp.Tool{Name: "godComment", Description: "生成网友的神评论"}, GodComment)
	// 73. 开车搞笑经历
	mcp.AddTool[DrivingStoryInput, map[string]interface{}](server, &mcp.Tool{Name: "drivingStory", Description: "生成开车时的搞笑经历"}, DrivingStory)
	// 74. 做饭翻车现场
	mcp.AddTool[CookingFailInput, map[string]interface{}](server, &mcp.Tool{Name: "cookingFail", Description: "生成做饭时的翻车现场"}, CookingFail)
	// 75. 自拍搞笑姿势
	mcp.AddTool[SelfiePoseInput, map[string]interface{}](server, &mcp.Tool{Name: "selfiePose", Description: "生成自拍时的搞笑姿势"}, SelfiePose)
	// 76. 聚会游戏推荐
	mcp.AddTool[PartyGameInput, map[string]interface{}](server, &mcp.Tool{Name: "partyGame", Description: "推荐聚会时的游戏"}, PartyGame)
	// 77. KTV必点歌曲
	mcp.AddTool[KTVSongInput, map[string]interface{}](server, &mcp.Tool{Name: "ktvSong", Description: "推荐KTV必点歌曲"}, KTVSong)
	// 78. 健身摸鱼行为
	mcp.AddTool[FitnessSlackInput, map[string]interface{}](server, &mcp.Tool{Name: "fitnessSlack", Description: "生成健身时的摸鱼行为"}, FitnessSlack)
	// 79. 网购省钱技巧
	mcp.AddTool[ShoppingSaveInput, map[string]interface{}](server, &mcp.Tool{Name: "shoppingSave", Description: "生成网购时的省钱技巧"}, ShoppingSave)
	// 80. 职场潜规则
	mcp.AddTool[WorkplaceRuleInput, map[string]interface{}](server, &mcp.Tool{Name: "workplaceRule", Description: "生成职场中的潜规则"}, WorkplaceRule)
	// 81. 校园奇葩规定
	mcp.AddTool[SchoolRuleInput, map[string]interface{}](server, &mcp.Tool{Name: "schoolRule", Description: "生成校园里的奇葩规定"}, SchoolRule)
	// 82. 恋爱甜蜜小事
	mcp.AddTool[LoveSweetInput, map[string]interface{}](server, &mcp.Tool{Name: "loveSweet", Description: "生成恋爱中的甜蜜小事"}, LoveSweet)
	// 83. 朋友暖心瞬间
	mcp.AddTool[FriendWarmInput, map[string]interface{}](server, &mcp.Tool{Name: "friendWarm", Description: "生成朋友间的暖心瞬间"}, FriendWarm)
	// 84. 家人温馨时刻
	mcp.AddTool[FamilyWarmInput, map[string]interface{}](server, &mcp.Tool{Name: "familyWarm", Description: "生成家人间的温馨时刻"}, FamilyWarm)
	// 85. 旅行暖心经历
	mcp.AddTool[TravelWarmInput, map[string]interface{}](server, &mcp.Tool{Name: "travelWarm", Description: "生成旅行中的暖心经历"}, TravelWarm)
	// 86. 生活小确幸
	mcp.AddTool[LittleHappinessInput, map[string]interface{}](server, &mcp.Tool{Name: "littleHappiness", Description: "生成生活中的小确幸"}, LittleHappiness)
	// 87. 动物可爱瞬间
	mcp.AddTool[AnimalCuteInput, map[string]interface{}](server, &mcp.Tool{Name: "animalCute", Description: "生成动物的可爱瞬间"}, AnimalCute)
	// 88. 自然美丽景色
	mcp.AddTool[NatureBeautyInput, map[string]interface{}](server, &mcp.Tool{Name: "natureBeauty", Description: "生成自然的美丽景色"}, NatureBeauty)
	// 89. 城市夜景
	mcp.AddTool[CityNightInput, map[string]interface{}](server, &mcp.Tool{Name: "cityNight", Description: "生成城市的夜景"}, CityNight)

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	log.Printf("MCP handler listening at %s", "http://localhost:8001")
	_ = http.ListenAndServe(":8003", handler)
	select {}
}
