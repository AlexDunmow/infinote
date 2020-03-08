export const objectSize = (obj: any): number => {
	let size = 0
	let key: string
	for (key in obj) {
		if (obj.hasOwnProperty(key)) size++
	}
	return size
}
