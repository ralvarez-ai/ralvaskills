import { Device } from '../components/Device';
import { Badge } from '../components/ui';
import { profileNames } from '../data/facts';

export function Writes() {
  return (
    <section className="k-sec" id="writes" aria-labelledby="d-wr-h" style={{ background: 'var(--bg2)' }}>
      <div className="wrap d-writes">
        <div className="rv">
          <div className="k-rule" aria-hidden="true" />
          <p className="d-call">
            <span className="c-no">3</span> after call 3
          </p>
          <h2 id="d-wr-h">It lands in two files you commit.</h2>
          <p className="k-lede">
            <code>install_skills</code> writes the skills into <code>.rsk/rsk.mod</code> and their exact versions into{' '}
            <code>.rsk/rsk.lock</code>. Whoever clones the repo runs <code>rsk install</code> and gets the same three, at
            the same versions. Think <code>go.mod</code> and <code>go.sum</code>, for your agent’s skills.
          </p>
        </div>
        <Device
          names={profileNames}
          second="install_skills · approved by you"
          agent
          caption={
            <>
              <Badge kind="sa">sample project</Badge> Same scratch folder as above. The <code>rsk new</code> lines are a
              real run (paths shortened); the ◆ line stands for call 3. The file contents are what{' '}
              <code>internal/install/apply.go</code> writes for those three skills, with versions from the registry on
              2026-11-04 — not a capture.
            </>
          }
        />
      </div>
    </section>
  );
}
