import typescript from '@rollup/plugin-typescript';
import { babel } from '@rollup/plugin-babel';
import { nodeResolve } from '@rollup/plugin-node-resolve';
import replace from '@rollup/plugin-replace';
import commonjs from '@rollup/plugin-commonjs';
import scss from 'rollup-plugin-scss';

export default {
    input: 'src/index.tsx',
    output: [
        {
            dir: './dist',
            format: 'cjs',
            sourcemap: true
        }, 
    ],
    plugins: [
        replace({
            'process.env.NODE_ENV': JSON.stringify('production')
        }),
        nodeResolve({ browser: true }),
        commonjs(),
        babel({
            extensions: ['.jsx', '.tsx'],
            babelHelpers: 'bundled',
            presets: ['@babel/preset-react'],
            exclude: 'node_modules/**'
        }),
        typescript(),
        scss({ output: './dist/index.css' }),
    ],
}