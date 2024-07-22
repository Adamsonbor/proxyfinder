import { useTheme } from '@mui/material/styles';
import Switch, { SwitchProps } from '@mui/material/Switch';

interface Props extends SwitchProps {
}

export default function ThemeSwitch({
	sx = {},
	...props
}: Props) {
	const theme = useTheme();

	return (
		<Switch
			focusVisibleClassName=".Mui-focusVisible"
			disableRipple
			sx={{
				...sx,
				width: 42,
				height: 26,
				padding: 0,
				'& .MuiButtonBase-root': {
					height: '100%',
				},
				'& .MuiSwitch-switchBase': {
					padding: '0px 4px',
					transitionDuration: '300ms',
					'&.Mui-checked': {
						transform: 'translateX(16px)',
						'& + .MuiSwitch-track': {
							backgroundColor: 'transparent',
							opacity: 0.5,
						},
						'& .MuiSwitch-thumb': {
							backgroundColor: theme.palette.blue,
						},

						'& .MuiSwitch-thumb::after': {
							backgroundImage: 'url("time-of-night-icon.svg")',
						},
					},
				},
				'& .MuiSwitch-thumb': {
					backgroundColor: theme.palette.lightBlue,
					boxSizing: 'border-box',
					width: 16,
					height: 16,
				},
				'& .MuiSwitch-thumb::after': {
					content: "''",
					display: 'block',
					backgroundImage: 'url("time-of-day-icon.svg")',
					backgroundRepeat: 'no-repeat',
					backgroundPosition: 'center',
					width: 16,
					height: 16,
				},
				'& .MuiSwitch-track': {
					borderRadius: 26 / 2,
					border: `2px solid ${theme.palette.lightGray}`,
					backgroundColor: 'transparent',
					opacity: 1,
					transition: theme.transitions.create(['background-color'], {
						duration: 500,
					}),
				},
			}}
			{...props}
		/>
	);
}
