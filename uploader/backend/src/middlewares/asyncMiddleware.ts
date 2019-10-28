export default fn =>
  (req, res, next)=> Promise.resolve(fn(req, res, next)).catch(e=> res.send({ok: 0, err: e}))