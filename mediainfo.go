package mediainfo

import (
	"encoding/json"
	"flag"
	"fmt"
	"os/exec"
	"strings"
)

var mediainfoBinary = flag.String("mediainfo-bin", "mediainfo", "the path to the mediainfo binary if it is not in the system $PATH")

type MediaInfo struct {
	Media	media	`json:"media"`
}

type media struct {
	ReferenceToFile	string	`json:"@ref"`
	Tracks		[]track	`json:"track"`
}

type track struct {
	Type                    string  `json:"@type"`
        ImageCount              string  `json:"ImageCount,omitempty"`
        FileExtension           string  `json:"FileExtension,omitempty"`
        Format                  string  `json:"Format,omitempty"`
        FileSize                string  `json:"FileSize,omitempty"`
        StreamSize              string  `json:"StreamSize,omitempty"`
        FileModifiedDate        string  `json:"File_Modified_Date,omitempty"`
        FileModifiedDateLocal   string  `json:"File_Modified_Date_Local,omitempty"`
        Width                   string  `json:"Width,omitempty"`
        Height                  string  `json:"Height,omitempty"`
        ColorSpace              string  `json:"ColorSpace,omitempty"`
        ChromaSubsampling       string  `json:"ChromaSubsampling,omitempty"`
        BitDepth                string  `json:"BitDepth,omitempty"`
        CompressionMode         string  `json:"Compression_Mode,omitempty"`
        Duration                string  `json:"Duration,omitempty"`
        OverallBitRateMode      string  `json:"OverallBitRateMode,omitempty"`
        OverallBitRate          string  `json:"OverallBitRate,omitempty"`
        FileName                string  `json:"FileName,omitempty"`
        FrameRate               string  `json:"FrameRate,omitempty"`
        WritingApplication      string  `json:"WritingApplication,omitempty"`
        BitRate                 string  `json:"BitRate,omitempty"`
        ScanType                string  `json:"ScanType,omitempty"`
        Interlacement           string  `json:"Interlacement,omitempty"`
        WritingLibrary          string  `json:"WritingLibrary,omitempty"`
        Channels                string  `json:"Channels,omitempty"`
        FormatInfo              string  `json:"FormatInfo,omitempty"`
        SamplingRate            string  `json:"SamplingRate,omitempty"`
        FormatProfile           string  `json:"FormatProfile,omitempty"`
}

func IsInstalled() bool {
	cmd := exec.Command(*mediainfoBinary)
	err := cmd.Run()
	if err != nil {
		if strings.HasSuffix(err.Error(), "no such file or directory") ||
			strings.HasSuffix(err.Error(), "executable file not found in %PATH%") ||
			strings.HasSuffix(err.Error(), "executable file not found in $PATH") {
			return false
		} else if strings.HasPrefix(err.Error(), "exit status 255") {
			return true
		}
	}
	return true
}

func (info MediaInfo) IsMedia() bool {
	return info.Media.Tracks[0].Type == "Video" || info.Media.Tracks[0].Type == "Audio"
}

func GetMediaInfo(fname string) ([]MediaInfo, error) {
	info := MediaInfo{}

	if !IsInstalled() {
		return []MediaInfo{info}, fmt.Errorf("Must install mediainfo")
	}
	out, err := exec.Command(*mediainfoBinary, "--Output=JSON", "-f", fname).Output()

	if err != nil {
		return []MediaInfo{info}, err
	}

	if out[0] == '[' {
		infoSlice := []MediaInfo{}
		if err := json.Unmarshal(out, &infoSlice); err != nil {
			return infoSlice, err
		}
		return infoSlice, nil
	}
	if err := json.Unmarshal(out, &info); err != nil {
		return []MediaInfo{info}, err
	}

	return []MediaInfo{info}, nil
}
