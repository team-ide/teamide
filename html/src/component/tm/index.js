import teamide from './base/index.js';

import './style/base.css';
import './style/layout.css';


if (typeof window !== 'undefined') {
  window.teamide = teamide;
  teamide.style.init();
}

export default teamide;