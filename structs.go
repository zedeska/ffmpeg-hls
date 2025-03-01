package main

type ffprobe_audio struct {
	Programs     []interface{} `json:"programs"`
	StreamGroups []interface{} `json:"stream_groups"`
	Streams      []struct {
		Index          int    `json:"index"`
		CodecName      string `json:"codec_name"`
		CodecLongName  string `json:"codec_long_name"`
		CodecType      string `json:"codec_type"`
		CodecTagString string `json:"codec_tag_string"`
		CodecTag       string `json:"codec_tag"`
		SampleFmt      string `json:"sample_fmt"`
		SampleRate     string `json:"sample_rate"`
		Channels       int    `json:"channels"`
		ChannelLayout  string `json:"channel_layout"`
		BitsPerSample  int    `json:"bits_per_sample"`
		InitialPadding int    `json:"initial_padding"`
		RFrameRate     string `json:"r_frame_rate"`
		AvgFrameRate   string `json:"avg_frame_rate"`
		TimeBase       string `json:"time_base"`
		StartPts       int    `json:"start_pts"`
		StartTime      string `json:"start_time"`
		BitRate        string `json:"bit_rate"`
		Disposition    struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
			NonDiegetic     int `json:"non_diegetic"`
			Captions        int `json:"captions"`
			Descriptions    int `json:"descriptions"`
			Metadata        int `json:"metadata"`
			Dependent       int `json:"dependent"`
			StillImage      int `json:"still_image"`
		} `json:"disposition"`
		Tags struct {
			Language                    string `json:"language"`
			BPSEng                      string `json:"BPS-eng"`
			DURATIONEng                 string `json:"DURATION-eng"`
			NUMBEROFFRAMESEng           string `json:"NUMBER_OF_FRAMES-eng"`
			NUMBEROFBYTESEng            string `json:"NUMBER_OF_BYTES-eng"`
			STATISTICSWRITINGAPPEng     string `json:"_STATISTICS_WRITING_APP-eng"`
			STATISTICSWRITINGDATEUTCEng string `json:"_STATISTICS_WRITING_DATE_UTC-eng"`
			STATISTICSTAGSEng           string `json:"_STATISTICS_TAGS-eng"`
		} `json:"tags"`
	} `json:"streams"`
}

type resolution struct {
	width   int
	height  int
	bitrate string
}
