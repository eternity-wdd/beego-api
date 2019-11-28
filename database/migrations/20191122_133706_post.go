package main

import (
	"github.com/astaxie/beego/migration"
)

// DO NOT MODIFY
type Post_20191122_133706 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &Post_20191122_133706{}
	m.Created = "20191122_133706"

	migration.Register("Post_20191122_133706", m)
}

// Run the migrations
func (m *Post_20191122_133706) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE post(`id` int(11) NOT NULL AUTO_INCREMENT,`title` varchar(128) NOT NULL,`body` longtext  NOT NULL,PRIMARY KEY (`id`))")
}

// Reverse the migrations
func (m *Post_20191122_133706) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `post`")
}
