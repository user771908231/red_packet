package api

import "casino_majianagv2/core/data"

type MjUser interface {
	GetStatus() *data.MjUserStatus
	Ready()
}
