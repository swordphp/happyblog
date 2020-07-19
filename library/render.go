package library

type Render struct{
    ErrNo  int         `json:"errNo"`
    ErrMsg string      `json:"errMsg"`
    Data   interface{} `json:"data"`
}


func (render *Render) Init() Render{
    render.ErrNo = 0
    render.ErrMsg = "success"
    return *render
}

/**
 * 设置错误信息
 *
 * param: int    errNo
 * param: string errMsg
 * return: Render
 */
func (render *Render) SetErr(errNo int,errMsg string) {
    render.ErrNo = errNo
    render.ErrMsg = errMsg
}

func (render *Render) SetData(data interface{}) {
    render.Data = data
}
