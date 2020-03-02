export interface ServerPlan {
	id: string // db uuid
	short_id: number // short so can say
	name: string // plan description
	regions: string[] // location in text
	price: number
	currency: string // short currency character
	memory: number // mb
	storage: number // gb
	storage_type: string // ssd, hdd
	vcpu: number
	bandwidth: number // gb
	deprecated: boolean
}

// refer to initial migration sql file
export interface ServerInstance extends ServerPlan {
	id: string // db uuid
	user_id: string
	plan_id: number
	name: string
	note: string // customer note
	bandwidth_allowed: number // gb
	bandwidth_used: number // gb
	location: string
	charge: number
	currency: string
	status: string // pending, active
	sys_online: boolean // system is switch on
	vcpu: number
	memory: number // mb
	storage: number // gb
	storage_type: string // ssd, hdd
	os: string // operating system and version text
	kvm_url: string
	ip_v4_address: string // ip v4 address
	ip_v4_mask: string
	ip_v4_gateway: string // gateway ip v4 address
	ip_v6_address: string // ip v6 address
	ip_v6_mask: string
	ip_v6_size: number
	created: string // iso 8601
}
