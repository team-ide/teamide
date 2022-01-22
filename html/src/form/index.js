import {  validator } from '@/form/base.js';
import form from '@/form/form.js';

form.build = validator.buildFormValidator;

export default form;