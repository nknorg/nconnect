export function assignDefined(target, ...sources) {
  for (let source of sources) {
    if (source) {
      for (let key of Object.keys(source)) {
        if (source[key] !== undefined) {
          target[key] = source[key];
        }
      }
    }
  }
  return target;
}
