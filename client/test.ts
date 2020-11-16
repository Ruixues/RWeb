import {rweb} from './rweb.ts'

const engine: rweb = new rweb("ws://127.0.0.1:1111/t");
engine.connect();