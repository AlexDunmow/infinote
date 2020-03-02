import * as React from "react"
import { useStyletron } from "baseui"

interface IProps {
	children: JSX.Element[]
	current: number
}

export const ProgressSteps = (props: IProps) => {
	const { children, current } = props

	const [useCss, theme] = useStyletron()
	const container = useCss({
		display: "flex",
		alignItems: "center",
		flexDirection: "column",
	})
	const progressBar = useCss({
		width: "100%",
		display: "flex",
		alignItems: "center",
		justifyContent: "space-between",
		marginBottom: "1em",
	})

	const content = useCss({
		width: "100%",
	})

	const progressBarElements: JSX.Element[] = []
	for (let i = 0; i < children.length; i++) {
		if (i != 0) {
			progressBarElements.push(<Tail key={"tail-" + i} isFilled={i <= current} />)
		}
		progressBarElements.push(<Circle key={"circle" + i} number={i + 1} isFilled={i <= current} />)
	}

	return (
		<div className={container}>
			<div className={progressBar}>{progressBarElements}</div>
			<div className={content}>{children[current]}</div>
		</div>
	)
}

const Circle = (props: { isFilled: boolean; number: number; key: string }) => {
	const { isFilled, number } = props
	const [useCss, theme] = useStyletron()
	const circle = useCss({
		width: "2rem",
		height: "2rem",
		display: "flex",
		alignItems: "center",
		justifyContent: "center",
		color: isFilled ? theme.colors.white : theme.colors.accent,
		borderRadius: "50%",
		flexShrink: 0,
		transitionDelay: isFilled ? "300ms" : "0ms",
		backgroundColor: isFilled ? theme.colors.progressStepsCompletedFill : theme.colors.accent100,
		zIndex: 20,
	})
	return (
		<div className={circle}>
			<strong>{number}</strong>
		</div>
	)
}

const Tail = (props: { isFilled: boolean; key: string }) => {
	const { isFilled } = props
	const [useCss, theme] = useStyletron()
	const tail = useCss({
		height: "0.5rem",
		margin: "0 -1px 0 -1px",
		flexGrow: 1,
		transitionTimingFunction: "linear",
		transitionProperty: "all",
		transitionDuration: "300ms",
		backgroundImage: `linear-gradient(to right, ${theme.colors.progressStepsCompletedFill} 50%, ${theme.colors.accent100} 50%)`,
		backgroundSize: "210% 100%",
		backgroundPosition: isFilled ? "left bottom" : "right bottom",
		zIndex: 10,
	})
	return <div className={tail} />
}
