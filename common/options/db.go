//
//  db.go
//  coreruntime_extenstions
//
//  Created by karim-w on 12/07/2025.
//

package options

import "github.com/karim-w/gopts"

type SQL_Database struct {
	MaxIdleConns   gopts.Option[int]
	MaxOpenConns   gopts.Option[int]
	PanicablePings bool
}

type Redis struct {
	PanicablePings bool
}
