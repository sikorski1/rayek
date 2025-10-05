import { motion } from "framer-motion";
import { Info } from "lucide-react";
import { ComponentProps } from "react";
import { Tooltip as ToolTipComponent } from "react-tooltip";
import styles from "./tooltip.module.scss";

type ReactTooltipProps = ComponentProps<typeof ToolTipComponent>;

type ToolTipProps = {
	name: string;
	index?: number;
	toolTipText: string;
} & ReactTooltipProps;

export default function ToolTip({ name, index = 0, toolTipText, ...props }: ToolTipProps) {
	return (
		<>
			<motion.div
				id={name + "toolTip"}
				className={styles.toolTipBox}
				initial={{ opacity: 0, y: 20, scale: 0.9 }}
				animate={{ opacity: 1, y: 0, scale: 1 }}
				transition={{
					duration: 0.4,
					ease: "easeOut",
					delay: index * 0.1,
				}}>
				<Info className={styles.toolTipIcon} />
			</motion.div>

			<ToolTipComponent
				role="tooltip"
				className={styles.toolTipTextBox}
				classNameArrow={styles.toolTipTextArrow}
				anchorSelect={"#" + name + "toolTip"}
				opacity={1}
				{...props}>
				{toolTipText}
			</ToolTipComponent>
		</>
	);
}
