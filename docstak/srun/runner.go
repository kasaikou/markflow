/*
Copyright 2024 Kasai Kou

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package srun

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"time"
)

type ScriptRunner struct {
	cmd *exec.Cmd
}

func NewScriptRunner(execPath string, script string) *ScriptRunner {
	runner := &ScriptRunner{
		cmd: exec.Command(execPath),
	}

	runner.cmd.Stdin = bytes.NewBufferString(script)
	return runner
}

func (sr *ScriptRunner) SetWorkingDir(dir string)   { sr.cmd.Dir = dir }
func (sr *ScriptRunner) SetEnviron(environ string)  { sr.cmd.Env = append(sr.cmd.Env, environ) }
func (sr *ScriptRunner) SetEnv(key, value string)   { sr.cmd.Env = append(sr.cmd.Env, key+"="+value) }
func (sr *ScriptRunner) Stdout() (io.Reader, error) { return sr.cmd.StdoutPipe() }
func (sr *ScriptRunner) Stderr() (io.Reader, error) { return sr.cmd.StderrPipe() }

func (sr *ScriptRunner) RunContext(ctx context.Context) (int, error) {

	var cmdErr error
	onFin := make(chan struct{}, 1)
	go func() {
		defer close(onFin)
		if cmdErr = sr.cmd.Start(); cmdErr != nil {
			return
		} else if cmdErr = sr.cmd.Wait(); cmdErr != nil {
			return
		}
	}()

	select {
	case <-ctx.Done():
		if err := sr.cmd.Process.Signal(os.Interrupt); err == nil {
			timer := time.NewTimer(10 * time.Second)
			select {
			case <-onFin:
				return sr.cmd.ProcessState.ExitCode(), cmdErr
			case <-timer.C:
			}
		}

		if err := sr.cmd.Process.Kill(); err != nil {
			return sr.cmd.ProcessState.ExitCode(), err
		}

		<-onFin
		return sr.cmd.ProcessState.ExitCode(), cmdErr

	case <-onFin:
		return sr.cmd.ProcessState.ExitCode(), cmdErr
	}
}
