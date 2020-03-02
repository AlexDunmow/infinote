import * as React from "react"
import { styled, useStyletron } from "baseui"
import { H6 } from "baseui/typography"

const LoadingSVG = require("../assets/images/loading.svg")

const LoadingImage = styled("img", {
	maxWidth: "100%"
})

export const Loading = () => {
	return <LoadingImage alt="loading" src={LoadingSVG} />
}
