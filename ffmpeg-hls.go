package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

func main() {

	file := flag.String("f", "null", "chemin d'accèes à une vidéo")
	folder := flag.String("d", "null", "chemin d'accèes à un dossier contenant des vidéos")

	flag.Parse()

	if *file == "null" && *folder == "null" {
		fmt.Println("Il faut spécifier une vidéo ou un dossier. utiliser -h pour plus d'info")
		return
	} else if *file != "null" && *folder != "null" {
		fmt.Println("Vous ne pouvez pas spécifier une vidéo et un dossier à la fois")
		return
	} else if *file != "null" {
		temp := strings.Split(*file, ".")
		output := strings.Join(temp[:len(temp)-1], ".")
		os.Mkdir(output, os.ModePerm)
		encode(*file, output)
	} else if *folder != "null" {
		entries, _ := os.ReadDir(*folder)
		for _, e := range entries {
			temp := strings.Split(*folder+"\\"+e.Name(), ".")
			if temp[len(temp)-1] == "mkv" || temp[len(temp)-1] == "mp4" {
				output := strings.Join(temp[:len(temp)-1], ".")
				os.Mkdir(output, os.ModePerm)
				encode(*folder+"\\"+e.Name(), output)
			}
		}
	}
}

func encode(input_file string, output_file string) {
	var result_audio ffprobe_audio
	var result_subtitle ffprobe_subtitle

	var resolution []resolution = []resolution{
		{
			width:   1920,
			height:  1080,
			bitrate: "2M",
		},
		{
			width:   1280,
			height:  720,
			bitrate: "1M",
		},
	}

	ffprobe_audio := exec.Command("ffprobe", "-v", "error", "-select_streams", "a", "-show_entries", "stream", "-of", "json", input_file)
	ffprobe_subtitle := exec.Command("ffprobe", "-v", "error", "-select_streams", "s", "-show_entries", "stream", "-of", "json", input_file)

	output_audio, _ := ffprobe_audio.Output()
	output_subtitle, _ := ffprobe_subtitle.Output()

	json.Unmarshal(output_audio, &result_audio)
	json.Unmarshal(output_subtitle, &result_subtitle)

	num_subs := len(result_subtitle.Streams)
	num_audio := len(result_audio.Streams)
	num_resolution := len(resolution)

	if num_subs > 0 {
		var subs []string
		for i := 0; i < num_subs; i++ {
			if result_subtitle.Streams[i].Tags.Title == "" {
				subs = append(subs, result_subtitle.Streams[i].Tags.Language)
			} else {
				subs = append(subs, strings.Replace(result_subtitle.Streams[i].Tags.Title, "/", "", -1))
			}
		}

		check_lang_dup(subs)

		for i := 0; i < num_subs; i++ {
			f := exec.Command("ffmpeg", "-i", input_file, "-map", "0:s:"+strconv.Itoa(i), "-f", "webvtt", output_file+"/"+subs[i]+".vtt")
			f.Stdout = os.Stdout
			f.Stderr = os.Stderr
			f.Run()
		}

	}

	var lang []string
	for i := 0; i < num_audio; i++ {
		lang = append(lang, result_audio.Streams[i].Tags.Language)
	}

	check_lang_dup(lang)

	var hwaccel string
	var encoder string
	gpu := getGPU()
	if strings.Contains(gpu, "AMD") {
		hwaccel = "d3d11va"
		encoder = "h264_amf"
	} else if strings.Contains(gpu, "NVIDIA") {
		hwaccel = "cuda"
		encoder = "h264_nvenc"
	}

	var ffmpeg_command []string = []string{"-hwaccel", hwaccel, "-i", input_file, "-filter_complex"}

	var filter_complex string
	filter_complex += fmt.Sprintf("[0:v:0]split=%d", num_resolution)

	for i := 0; i < num_resolution; i++ {
		filter_complex += fmt.Sprintf("[v%d]", i)
	}

	for i := 0; i < num_resolution; i++ {
		filter_complex += fmt.Sprintf(";[v%d]scale=%d:%d[v%dout]", i, resolution[i].width, resolution[i].height, i)
	}

	ffmpeg_command = append(ffmpeg_command, filter_complex)

	for i := 0; i < num_resolution; i++ {
		ffmpeg_command = append(ffmpeg_command, []string{"-map", fmt.Sprintf("[v%dout]", i), fmt.Sprintf("-c:v:%d", i), encoder, fmt.Sprintf("-b:v:%d", i), resolution[i].bitrate, "-preset", "medium", "-profile:v", "main", "-pix_fmt", "yuv420p", fmt.Sprintf("-s:v:%d", i), fmt.Sprintf("%dx%d", resolution[i].width, resolution[i].height)}...)
	}

	for i := 0; i < num_audio; i++ {
		ffmpeg_command = append(ffmpeg_command, []string{"-map", fmt.Sprintf("0:a:%d", i), fmt.Sprintf("-c:a:%d", i), "aac", "-ac", "2", fmt.Sprintf("-b:a:%d", i), "128k", fmt.Sprintf("-metadata:s:a:%d", i), fmt.Sprintf("language=%s", lang[i])}...)
	}

	ffmpeg_command = append(ffmpeg_command, []string{"-f", "hls", "-hls_time", "10", "-hls_list_size", "0", "-hls_flags", "independent_segments", "-master_pl_name", "master.m3u8", "-var_stream_map"}...)

	var stream_map string

	for i := 0; i < num_resolution; i++ {
		stream_map += fmt.Sprintf("v:%d,agroup:audio,name:%dp ", i, resolution[i].height)
	}

	for i := 0; i < num_audio; i++ {
		if i == num_audio-1 {
			stream_map += fmt.Sprintf("a:%d,agroup:audio,language:%s,name:audio_%s,default:yes", i, strings.ToUpper(lang[i]), lang[i])
			break
		}
		stream_map += fmt.Sprintf("a:%d,agroup:audio,language:%s,name:audio_%s ", i, strings.ToUpper(lang[i]), lang[i])
	}

	ffmpeg_command = append(ffmpeg_command, stream_map)
	ffmpeg_command = append(ffmpeg_command, []string{"-hls_segment_filename", output_file + "/%v/segment_%03d.ts", output_file + "/%v/manifest.m3u8"}...)

	//fmt.Println(ffmpeg_command)
	cmd := exec.Command("ffmpeg", ffmpeg_command...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("success")

}

func check_lang_dup(list []string) {
	temp := list
	var count []lang_count
	var found bool = false

	for i := 0; i < len(list); i++ {
		if len(temp) <= 1 {
			break
		}
		x := temp[0]
		temp = temp[1:]

		if slices.Contains(temp, x) {
			found = false
			for ii := 0; i < len(count); ii++ {
				if count[ii].lang == x {
					found = true
					count[ii].count += 1
					list[i] = x + strconv.Itoa(count[ii].count)
					break
				}
			}

			if !found {
				count = append(count, lang_count{lang: x, count: 2})
				list[i] = x + strconv.Itoa(2)
			}
		}
	}
}

func getGPU() string {
	Info := exec.Command("cmd", "/C", "wmic path win32_VideoController get name")
	History, _ := Info.Output()

	return strings.TrimSpace(strings.Replace(string(History), "Name", "", -1))
}
