/*
 * @Description:
 * @Author: yujiajie
 * @Date: 2024-03-16 22:45:15
 * @LastEditTime: 2024-03-16 22:45:51
 * @LastEditors: yujiajie
 */
package config

type LoggerConfig struct {
	LogLevel   string `json:"level"`
	LogPath    string `json:"logpath"`
	MaxSize    int    `json:"maxsize"`
	MaxAge     int    `json:"age"`
	MaxBackups int    `json:"backups"`
	Compress   string `json:"compress"`
}
