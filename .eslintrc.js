module.exports = {
    "env": {
        "browser": true,
        "es6": true
    },
    "extends": "google",
    "globals": {
        "Atomics": "readonly",
        "SharedArrayBuffer": "readonly"
    },
    "parserOptions": {
        "ecmaVersion": 2018,
        "sourceType": "module"
    },
    "rules": {
      'max-len': [2, 120, 2],
    'object-curly-spacing': ['error', 'always'],
    }
};